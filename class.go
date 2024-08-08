package mruby

import (
	"fmt"
)

var (
	FlagUndefinedAllocate = uint32(1 << 6)
	FlagClassIsInherited  = uint32(1 << 17)
	FlagClassIsOrigin     = uint32(1 << 18)
	FlagClassIsPrepended  = uint32(1 << 19)
)

type methodTable map[Symbol]Method

var _ RClass = &Class{}

type RClass interface {
	RObject
	Super() RClass
	mtPut(Symbol, Method)
	mtGet(Symbol) Method
}

type class struct {
	super RClass
	class RClass
	flags uint32
	mt    methodTable
	iv    ivTable
}

type Class struct {
	class
}

type SingletonClass struct {
	class
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

func (mrb *State) vmDefineClass(outer Value, super Value, id Symbol) (RClass, error) {
	var superClass RClass
	if super != nil {
		if ClassP(super) {
			superClass = super.(RClass)
		} else {
			panic("super is not a class")
		}
	}

	if !ClassPointerP(outer) {
		panic("outer is not a class or module")
	}
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

func (c *class) Class() RClass {
	return c.class
}

func (c *class) Flags() uint32 {
	return c.flags
}

func (c *class) ivPut(sym Symbol, val Value) {
	if c.iv == nil {
		c.iv = make(ivTable)
	}

	c.iv[sym] = val
}

func (c *class) ivGet(sym Symbol) Value {
	if c.iv == nil {
		return nil
	}

	return c.iv[sym]
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

	basicObject.class.class = classClass
	objectClass.class.class = classClass
	moduleClass.class.class = classClass
	classClass.class.class = classClass

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
