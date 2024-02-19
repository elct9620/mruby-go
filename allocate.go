package mruby

func (mrb *State) AllocClass() *Class {
	return &Class{
		class: class{
			object: object{
				class: mrb.ClassClass,
			},
		},
	}
}

func (mrb *State) AllocSingletonClass() *SingletonClass {
	return &SingletonClass{
		class: class{
			object: object{
				class: mrb.ClassClass,
			},
		},
	}
}

func (mrb *State) AllocModule() *Module {
	return &Module{
		class: class{
			object: object{
				class: mrb.ModuleClass,
			},
		},
	}
}

func (mrb *State) AllocObject(class RClass) RObject {
	return &Object{object{class}, nil}
}

func (mrb *State) AllocException(class RClass) RException {
	return &Exception{object{class}, ""}
}
