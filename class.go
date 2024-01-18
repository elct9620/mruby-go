package mruby

type methodTable map[Symbol]*Method
type mt = methodTable

var _ RClass = &Class{}

type RClass interface {
	RObject
	Super() RClass
	LookupMethod(Symbol) *Method
	Set(Symbol, Value)
	Get(Symbol) Value
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

func (mrb *State) NewClass(super *Class) *Class {
	return newClass(mrb, super)
}

func (mrb *State) DefineClass(outer Value, super Value, name string) *Class {
	return mrb.DefineClassById(outer, super, mrb.Intern(name))
}

func (mrb *State) DefineClassById(outer Value, super Value, id Symbol) *Class {
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
	class := newClass(mrb, superClass)

	// NOTE: setup_class()
	mrb.nameClass(class, outerModule, id)
	outerModule.Set(id, NewObjectValue(class))
	return class
}

func (mrb *State) ClassName(class RClass) string {
	if class == nil {
		return ""
	}

	name := class.Get(_classname(mrb))
	if name == nil {
		return ""
	}

	return name.(string)
}

func newClass(mrb *State, super *Class) *Class {
	class := &Class{
		class: class{
			mt: make(methodTable),
			iv: make(ivTable),
		},
	}

	if super != nil {
		class.super = super
	} else {
		class.super = mrb.ObjectClass
	}

	return class
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

	class.Set(nsym, name)
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
	basicObject := newClass(mrb, nil)
	objectClass := newClass(mrb, basicObject)
	mrb.ObjectClass = objectClass
	moduleClass := newClass(mrb, mrb.ObjectClass)
	mrb.ModuleClass = moduleClass
	classClass := newClass(mrb, mrb.ModuleClass)
	mrb.ClassClass = classClass

	basicObject.object.class = classClass
	objectClass.object.class = classClass
	moduleClass.object.class = classClass

	mrb.DefineConstById(basicObject, _BasicObject(mrb), NewObjectValue(basicObject))
	mrb.DefineConstById(objectClass, _Object(mrb), NewObjectValue(objectClass))
	mrb.DefineConstById(objectClass, _Module(mrb), NewObjectValue(moduleClass))
	mrb.DefineConstById(objectClass, _Class(mrb), NewObjectValue(classClass))

	mrb.nameClass(basicObject, nil, _BasicObject(mrb))
	mrb.nameClass(objectClass, nil, _Object(mrb))
	mrb.nameClass(moduleClass, nil, _Module(mrb))
	mrb.nameClass(classClass, nil, _Class(mrb))
}
