package mruby

var _ RClass = &Module{}

type Module struct {
	class
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
