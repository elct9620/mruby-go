package mruby

import "fmt"

type object struct {
	class RClass
}

type RObject interface {
	Class() RClass
}

var _ RObject = &Object{}

type Object struct {
	object
}

func (obj *object) Class() RClass {
	return nil
}

func initObject(mrb *State) {
	mrb.FalseClass = mrb.NewClass(mrb.ObjectClass)
	mrb.TrueClass = mrb.NewClass(mrb.ObjectClass)

	mrb.ObjectClass.DefineMethod(mrb, "new", &Method{
		Function: func(mrb *State, recv Value) Value {
			args := mrb.GetArgv()
			argc := mrb.GetArgc()

			super := mrb.ObjectClass
			if argc > 0 {
				super = args[0].(*Class)
			}

			return &Object{object{super}}
		},
	})

	mrb.ObjectClass.DefineMethod(mrb, "puts", &Method{
		Function: func(mrb *State, recv Value) Value {
			args := mrb.GetArgv()
			fmt.Println(args...)
			return args[0]
		},
	})
}
