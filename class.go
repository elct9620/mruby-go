package mruby

import (
	"fmt"
)

var (
	FlagUndefinedAllocate = uint32(1 << 6)
	FlagClassIsInherited  = uint32(1 << 17)
	FlagClassIsOrigin     = uint32(1 << 18)
	FlagClassIsPrepended  = uint32(1 << 19)
	FlagObjectIsFrozen    = uint32(1 << 20)
)

type methodTable map[Symbol]Method
type mt = methodTable

var _ RClass = &Class{}
var _ RClass = &SingletonClass{}

type RClass interface {
	RObject
	Super() RClass
	setClass(RClass)
	setSuper(RClass)
	getMethodTable() methodTable
	setMethodTable(methodTable)
	mtPut(Symbol, Method)
	mtGet(Symbol) Method
}

type class struct {
	Object
	super RClass
	mt
	iv
}

type Class struct {
	class
}

type SingletonClass struct {
	class
}

type InheritClass struct {
	class
}

func ClassPointerP(v Value) bool {
	switch v.(type) {
	case *Class, *SingletonClass, *Module:
		return true
	default:
		return false
	}
}

func (mrb *State) Class(v Value) RClass {
	switch v := v.(type) {
	case RObject:
		return v.Class()
	case []Value:
		return mrb.ArrayClass
	case bool:
		if v {
			return mrb.TrueClass
		}

		return mrb.FalseClass
	}

	return nil
}

func (mrb *State) DefineClassId(name Symbol, super RClass) (RClass, error) {
	return mrb.defineClass(name, super, mrb.ObjectClass)
}

func (mrb *State) ObjectClassName(obj Value) string {
	return mrb.ClassName(mrb.Class(obj))
}

func (mrb *State) ClassName(class RClass) string {
	if class == nil {
		return ""
	}

	name := mrb.ObjectInstanceVariableGet(class, _classname(mrb))
	if name == nil {
		return ""
	}

	return name.(string)
}

func (mrb *State) ClassNew(super RClass) (*Class, error) {
	if super != nil {
		mrb.checkInheritable(super)
	}

	class := mrb.bootDefineClass(super)
	if super != nil {
		class.flags |= super.Flags() & FlagUndefinedAllocate
	}

	err := mrb.prepareSingletonClass(class)
	if err != nil {
		return nil, err
	}

	return class, nil
}

func (mrb *State) DefineMethodId(class RClass, name Symbol, function Function) {
	method := newMethodFromFunc(function)
	mrb.defineMethodRaw(class, name, method)
}

func (mrb *State) VmFindMethod(recv Value, class RClass, mid Symbol) Method {
	c := class
	for c != nil {
		m := c.mtGet(mid)
		if m != nil {
			return m
		}

		c = c.Super()
	}

	return nil
}

func (mrb *State) checkIfClassOrModule(object Value) {
	if !ClassPointerP(object) {
		mrb.Raisef(nil, "%v is not a class/module", object)
	}
}

func (mrb *State) vmDefineClass(outer Value, super Value, id Symbol) (RClass, error) {
	var superClass RClass
	if super != nil {
		if ClassP(super) {
			superClass = super.(RClass)
		} else {
			mrb.Raisef(nil, "superclass must be a Class (%T given)", super)
		}
	}

	mrb.checkIfClassOrModule(outer)
	outerModule := outer.(RClass)

	return mrb.defineClass(id, superClass, outerModule)
}

func (mrb *State) defineClass(name Symbol, super RClass, outer RClass) (RClass, error) {
	class, err := mrb.ClassNew(super)
	if err != nil {
		return nil, err
	}

	mrb.setupClass(class, outer, name)
	return class, nil
}

func (mrb *State) setupClass(class RClass, outer RClass, id Symbol) {
	mrb.nameClass(class, outer, id)
	mrb.ObjectInstanceVariableSetForce(outer, id, NewObjectValue(class))
}

func (mrb *State) prepareSingletonClass(obj RObject) error {
	if obj.Class() == nil {
		return fmt.Errorf("class not found on object type %T", obj)
	}

	if _, ok := obj.Class().(*SingletonClass); ok {
		return nil
	}

	singletonClass := mrb.AllocSingletonClass()
	if class, ok := obj.(*Class); ok {
		if class.Super() != nil {
			singletonClass.super = class.Super()
		} else {
			singletonClass.super = mrb.ClassClass
		}
	} else if _, ok := obj.(*SingletonClass); ok {
		// NOTE: find origin class
	} else {
		singletonClass.super = obj.Class()
		err := mrb.prepareSingletonClass(singletonClass)
		if err != nil {
			return err
		}
	}

	mrb.ObjectInstanceVariableSetForce(singletonClass, _attached(mrb), obj)
	singletonClass.flags = singletonClass.Flags() & FlagObjectIsFrozen

	return nil
}

func (mrb *State) includeClassNew(module, super RClass) RClass {
	ic := mrb.AllocInheritClass()
	if _, isIClass := module.(*InheritClass); isIClass {
		module = module.Class()
	}

	module = findOrigin(module)
	ic.setMethodTable(module.getMethodTable())
	ic.setSuper(super)

	if _, isIClass := module.(*InheritClass); isIClass {
		ic.setClass(module.Class())
	} else {
		ic.setClass(module)
	}

	return ic
}

func (mrb *State) includeModuleAt(class, insPos, module RClass, searchSuper bool) int {
	for module != nil {
		parentClass := class.Super()

		if module.Flags()&FlagClassIsPrepended != 0 {
			goto Skip
		}

		for parentClass != nil {
			if !searchSuper {
				break
			}

			parentClass = parentClass.Super()
		}

		{
			ic := mrb.includeClassNew(module, insPos.Super())
			insPos.setSuper(ic)
			insPos = ic
		}
	Skip:
		module = module.Super()
	}

	return 0
}

func (mrb *State) IncludeModule(class, module RClass) error {
	mrb.includeModuleAt(class, findOrigin(class), module, true)

	return nil
}

func (mrb *State) bootDefineClass(super RClass) *Class {
	class := mrb.AllocClass()

	if super != nil {
		class.super = super
		class.flags |= FlagClassIsInherited
	} else {
		class.super = mrb.ObjectClass
	}

	class.mt = make(methodTable)
	return class
}

func (mrb *State) defineMethodRaw(class RClass, name Symbol, method Method) {
	class.mtPut(name, method)
}

func (c *class) setSuper(super RClass) {
	c.super = super
}

func (c *class) setClass(class RClass) {
	c.class = class
}

func (c *class) getMethodTable() methodTable {
	return c.mt
}

func (c *class) setMethodTable(mt methodTable) {
	c.mt = mt
}

func (c *class) ivPut(sym Symbol, val Value) {
	if c.iv == nil {
		c.iv = make(iv)
	}

	c.iv.Put(sym, val)
}

func (c *class) ivGet(sym Symbol) Value {
	if c.iv == nil {
		return nil
	}

	return c.iv.Get(sym)
}

func (c *class) mtPut(sym Symbol, method Method) {
	if c.mt == nil {
		c.mt = make(methodTable)
	}

	c.mt[sym] = method
}

func (c *class) mtGet(sym Symbol) Method {
	if c.mt == nil {
		return nil
	}

	return c.mt[sym]
}

func (c *class) Super() RClass {
	return c.super
}

func (mrb *State) nameClass(class RClass, outer RClass, id Symbol) {
	name := mrb.SymbolName(id)
	nsym := _classname(mrb)

	if mrb.ObjectInstanceVariableDefined(class, nsym) {
		return
	}

	mrb.ObjectInstanceVariableSetForce(class, nsym, name)
}

func (mrb *State) initClassNew(class RClass) {
	mrb.DefineMethodId(class, _new(mrb), allocObject)
}

func (mrb *State) checkInheritable(super RClass) {
	if _, ok := super.(*Class); !ok {
		mrb.Raisef(nil, "superclass must be a Class (%T given)", super)
	}

	if _, ok := super.(*SingletonClass); ok {
		mrb.Raisef(nil, "can't make subclass of singleton class")
	}

	if super == mrb.ClassClass {
		mrb.Raisef(nil, "can't make subclass of Class")
	}
}

func findOrigin(class RClass) RClass {
	if class == nil {
		return nil
	}

	ret := class
	if (ret.Flags() & FlagClassIsPrepended) != 0 {
		ret = ret.Super()

		for {
			if (ret.Flags() & FlagClassIsOrigin) == 0 {
				ret = ret.Super()
			}
		}
	}

	return ret
}

func allocObject(mrb *State, self Value) Value {
	args := mrb.GetArgv()
	argc := mrb.GetArgc()

	class := self.(RClass)
	if argc > 0 {
		class = args[0].(RClass)
	}

	return mrb.AllocObject(class)
}

func initClass(mrb *State) (err error) {
	basicObject := mrb.bootDefineClass(nil)
	objectClass := mrb.bootDefineClass(basicObject)
	mrb.ObjectClass = objectClass
	moduleClass := mrb.bootDefineClass(objectClass)
	mrb.ModuleClass = moduleClass
	classClass := mrb.bootDefineClass(moduleClass)
	mrb.ClassClass = classClass

	basicObject.Object.class = classClass
	objectClass.Object.class = classClass
	moduleClass.Object.class = classClass
	classClass.Object.class = classClass

	for _, class := range []RClass{basicObject, objectClass, moduleClass, classClass} {
		err = mrb.prepareSingletonClass(class)
		if err != nil {
			return
		}
	}

	mrb.DefineConstById(basicObject, _BasicObject(mrb), NewObjectValue(basicObject))
	mrb.DefineConstById(objectClass, _Object(mrb), NewObjectValue(objectClass))
	mrb.DefineConstById(objectClass, _Module(mrb), NewObjectValue(moduleClass))
	mrb.DefineConstById(objectClass, _Class(mrb), NewObjectValue(classClass))

	mrb.nameClass(basicObject, nil, _BasicObject(mrb))
	mrb.nameClass(objectClass, nil, _Object(mrb))
	mrb.nameClass(moduleClass, nil, _Module(mrb))
	mrb.nameClass(classClass, nil, _Class(mrb))

	mrb.initClassNew(classClass)
	mrb.topSelf = mrb.AllocObject(objectClass)

	return
}
