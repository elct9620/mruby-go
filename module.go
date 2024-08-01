package mruby

var _ RClass = &Module{}

type Module struct {
	super RClass
	class RClass
	flags uint32
	mt    methodTable
	iv    ivTable
}

func (c *Module) Class() RClass {
	return c.class
}

func (c *Module) Flags() uint32 {
	return c.flags
}

func (c *Module) ivPut(sym Symbol, val Value) {
	if c.iv == nil {
		c.iv = make(ivTable)
	}

	c.iv[sym] = val
}

func (c *Module) ivGet(sym Symbol) Value {
	if c.iv == nil {
		return nil
	}

	return c.iv[sym]
}

func (c *Module) mtPut(sym Symbol, method Method) {
	if c.mt == nil {
		c.mt = make(methodTable)
	}

	c.mt[sym] = method
}

func (c *Module) mtGet(sym Symbol) Method {
	if c.mt == nil {
		return nil
	}

	return c.mt[sym]
}

func (c *Module) Super() RClass {
	return c.super
}

func (mrb *State) DefineModuleId(name Symbol) RClass {
	return defineModule(mrb, name, mrb.ObjectClass)
}

func (mrb *State) vmDefineModule(outer RClass, name Symbol) RClass {
	return defineModule(mrb, name, outer)
}

func defineModule(mrb *State, name Symbol, outer RClass) *Module {
	module := mrb.AllocModule()
	module.mt = make(methodTable)
	mrb.setupClass(module, outer, name)

	return module
}
