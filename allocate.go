package mruby

func (mrb *State) AllocClass() *Class {
	return &Class{
		class: class{
			class: mrb.ClassClass,
		},
	}
}

func (mrb *State) AllocSingletonClass() *SingletonClass {
	return &SingletonClass{
		class: class{
			class: mrb.ClassClass,
		},
	}
}

func (mrb *State) AllocModule() *Module {
	return &Module{
		class: class{
			class: mrb.ModuleClass,
		},
	}
}

func (mrb *State) AllocObject(class RClass) RObject {
	return &Object{
		class: class,
		iv:    nil,
	}
}

func (mrb *State) AllocException(class RClass) RException {
	return &Exception{
		Object: Object{
			class: class,
		},
		message: "",
	}
}
