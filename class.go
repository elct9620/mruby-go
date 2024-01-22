package mruby

import "errors"

var (
	ErrObjectClassNotExists = errors.New("Object class not exists")
)

type methodTable map[Symbol]*Method
type mt = methodTable

var _ RClass = &Class{}

type RClass interface {
	RObject
	Super() RClass
	LookupMethod(Symbol) *Method
	DefineMethod(*State, string, *Method)
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

func (mrb *State) NewClass(super *Class) *Class {
	return mrb.bootDefineClass(super)
}

func (mrb *State) vmDefineClass(outer Value, super Value, id Symbol) *Class {
	superClass, ok := super.(*Class)
	if super != nil && !ok {
		panic("super is not a class")
	}

	// NOTE: check_if_class_or_module
	outerModule, ok := outer.(*Class)
	if !ok {
		panic("outer is not a class or module")
	}

	// NOTE: check constant defined in outer
	class := mrb.bootDefineClass(superClass)

	// NOTE: setup_class()
	mrb.nameClass(class, outerModule, id)
	mrb.ObjectInstanceVariableSetForce(outerModule, id, NewObjectValue(class))
	return class
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

func (mrb *State) bootDefineClass(super *Class) *Class {
	class := mrb.AllocClass()
	if super != nil {
		class.super = super
	} else {
		class.super = mrb.ObjectClass
	}

	class.mt = make(methodTable)
	return class
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

func (c *class) DefineMethod(mrb *State, name string, m *Method) {
	mid := mrb.Intern(name)
	c.mt[mid] = m
}

func (c *class) Super() RClass {
	return c.super
}

func (c *class) LookupMethod(mid Symbol) *Method {
	if c.mt[mid] != nil {
		return c.mt[mid]
	}

	super := c.super
	for super != nil {
		m := super.LookupMethod(mid)
		if m != nil {
			return m
		}

		super = super.Super()
	}

	return nil
}

func (mrb *State) nameClass(class RClass, outer RClass, id Symbol) {
	name := mrb.SymbolName(id)
	nsym := _classname(mrb)

	mrb.ObjectInstanceVariableSetForce(class, nsym, name)
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

func (mrb *State) FindMethod(recv Value, class RClass, mid Symbol) *Method {
	m := class.LookupMethod(mid)
	if m != nil {
		return m
	}

	return nil
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
