package mruby

var _ RClass = &Module{}

type Module struct {
	class
}

func (mrb *State) DefineModule(name string) *Module {
	return newModule(mrb)
}

func newModule(mrb *State) *Module {
	return &Module{
		class: class{
			super: mrb.ModuleClass,
			mt:    make(methodTable),
			iv:    make(ivTable),
		},
	}
}
