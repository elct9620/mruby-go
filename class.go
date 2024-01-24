package mruby

import (
	"errors"
)

var (
	ErrObjectClassNotExists = errors.New("Object class not exists")
)

type methodTable map[Symbol]Method
type mt = methodTable

var _ RClass = &Class{}

type RClass interface {
	RObject
	Super() RClass
	mtPut(Symbol, Method)
	mtGet(Symbol) Method
}

type class struct {
	object
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

func (mrb *State) ClassOf(v Value) *Class {
	switch v.(type) {
	case *Object:
		return mrb.ObjectClass
	case bool:
		if v == false {
			return mrb.FalseClass
		}

		return mrb.TrueClass
	}

	return nil
}

func (mrb *State) DefineClassId(name Symbol, super *Class) (*Class, error) {
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
	class := mrb.bootDefineClass(super)
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

func (mrb *State) vmDefineClass(outer Value, super Value, id Symbol) (*Class, error) {
	superClass, ok := super.(*Class)
	if super != nil && !ok {
		panic("super is not a class")
	}

	// NOTE: check_if_class_or_module
	outerModule, ok := outer.(*Class)
	if !ok {
		panic("outer is not a class or module")
	}

	return mrb.defineClass(id, superClass, outerModule)
}

func (mrb *State) defineClass(name Symbol, super RClass, outer RClass) (*Class, error) {
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
		return ErrObjectClassNotExists
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
	} else {
		class.super = mrb.ObjectClass
	}

	class.mt = make(methodTable)
	return class
}

func (mrb *State) defineMethodRaw(class RClass, name Symbol, method Method) {
	class.mtPut(name, method)
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

	mrb.ObjectInstanceVariableSetForce(class, nsym, name)
}

func initClass(mrb *State) {
	basicObject := mrb.bootDefineClass(nil)
	objectClass := mrb.bootDefineClass(basicObject)
	mrb.ObjectClass = objectClass
	moduleClass := mrb.bootDefineClass(objectClass)
	mrb.ModuleClass = moduleClass
	classClass := mrb.bootDefineClass(moduleClass)
	mrb.ClassClass = classClass

	basicObject.object.class = classClass
	objectClass.object.class = classClass
	moduleClass.object.class = classClass

	mrb.prepareSingletonClass(basicObject) // nolint: errcheck
	mrb.prepareSingletonClass(objectClass) // nolint: errcheck
	mrb.prepareSingletonClass(moduleClass) // nolint: errcheck
	mrb.prepareSingletonClass(classClass)  // nolint: errcheck

	mrb.DefineConstById(basicObject, _BasicObject(mrb), NewObjectValue(basicObject))
	mrb.DefineConstById(objectClass, _Object(mrb), NewObjectValue(objectClass))
	mrb.DefineConstById(objectClass, _Module(mrb), NewObjectValue(moduleClass))
	mrb.DefineConstById(objectClass, _Class(mrb), NewObjectValue(classClass))

	mrb.nameClass(basicObject, nil, _BasicObject(mrb))
	mrb.nameClass(objectClass, nil, _Object(mrb))
	mrb.nameClass(moduleClass, nil, _Module(mrb))
	mrb.nameClass(classClass, nil, _Class(mrb))
}
