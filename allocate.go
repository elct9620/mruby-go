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
