package mruby

type MethodTable map[string]*Method
type RClass struct {
	mt MethodTable
}

func NewClass() *RClass {
	return &RClass{
		mt: make(MethodTable),
	}
}

func (c *RClass) DefineMethod(name string, m *Method) {
	c.mt[name] = m
}

func (c *RClass) LookupMethod(name string) *Method {
	return c.mt[name]
}

func (mrb *State) ClassOf(v Value) *RClass {
	switch v.(type) {
	case *RObject:
		return mrb.objectClass
	case bool:
		if v == false {
			return mrb.falseClass
		}

		return mrb.trueClass
	}

	return nil
}

func (mrb *State) FindMethod(recv Value, class *RClass, mid string) *Method {
	m := class.LookupMethod(mid)
	if m != nil {
		return m
	}

	return nil
}
