package mruby

type MethodTable map[Symbol]*Method
type RClass struct {
	RBasic
	super *RClass
	mt    MethodTable
}

func (mrb *State) NewClass(super *RClass) *RClass {
	return newClass(mrb, super)
}

func (mrb *State) DefineModule(name string) *RClass {
	return newModule(mrb)
}

func newModule(mrb *State) *RClass {
	return &RClass{
		super: mrb.ModuleClass,
		mt:    make(MethodTable),
	}
}

func newClass(mrb *State, super *RClass) *RClass {
	class := &RClass{
		mt: make(MethodTable),
	}

	if super != nil {
		class.super = super
	}

	return class
}

func (c *RClass) DefineMethod(mrb *State, name string, m *Method) {
	mid := mrb.Intern(name)
	c.mt[mid] = m
}

func (c *RClass) LookupMethod(mid Symbol) *Method {
	return c.mt[mid]
}

func (mrb *State) ClassOf(v Value) *RClass {
	switch v.(type) {
	case *RObject:
		return mrb.ObjectClass
	case bool:
		if v == false {
			return mrb.FalseClass
		}

		return mrb.TrueClass
	}

	return nil
}

func (mrb *State) FindMethod(recv Value, class *RClass, mid Symbol) *Method {
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
}
