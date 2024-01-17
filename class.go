package mruby

type methodTable map[Symbol]*Method
type mt = methodTable

var _ Basic = &Class{}

type Class struct {
	object
	super *Class
	mt
	iv
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

func (mrb *State) ClassName(class *Class) string {
	if class == nil {
		return ""
	}

	name := class.Get(_classname(mrb))
	if name == nil {
		return ""
	}

	return name.(string)
}

func (mrb *State) DefineModule(name string) *Class {
	return newModule(mrb)
}

func newModule(mrb *State) *Class {
	return &Class{
		super: mrb.ModuleClass,
		mt:    make(methodTable),
		iv:    make(ivTable),
	}
}

func newClass(mrb *State, super *Class) *Class {
	class := &Class{
		mt: make(methodTable),
		iv: make(ivTable),
	}

	if super != nil {
		class.super = super
	} else {
		class.super = mrb.ObjectClass
	}

	return class
}

func (c *Class) DefineMethod(mrb *State, name string, m *Method) {
	mid := mrb.Intern(name)
	c.mt[mid] = m
}

func (c *Class) LookupMethod(mid Symbol) *Method {
	class := c

	for class != nil {
		m := class.mt[mid]
		if m != nil {
			return m
		}

		class = class.super
	}

	return nil
}

func (mrb *State) nameClass(class *Class, outer *Class, id Symbol) {
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

func (mrb *State) FindMethod(recv Value, class *Class, mid Symbol) *Method {
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

	basicObject.class = classClass
	objectClass.class = classClass
	moduleClass.class = classClass

	mrb.DefineConstById(basicObject, _BasicObject(mrb), NewObjectValue(basicObject))
	mrb.DefineConstById(objectClass, _Object(mrb), NewObjectValue(objectClass))
	mrb.DefineConstById(objectClass, _Module(mrb), NewObjectValue(moduleClass))
	mrb.DefineConstById(objectClass, _Class(mrb), NewObjectValue(classClass))

	mrb.nameClass(basicObject, nil, _BasicObject(mrb))
	mrb.nameClass(objectClass, nil, _Object(mrb))
	mrb.nameClass(moduleClass, nil, _Module(mrb))
	mrb.nameClass(classClass, nil, _Class(mrb))
}
