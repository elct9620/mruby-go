package mruby

func (mrb *State) allocSingletonClass() *SingletonClass {
	return &SingletonClass{
		class: class{
			object: object{
				class: mrb.ClassClass,
			},
		},
	}
}
