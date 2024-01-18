package mruby

import "fmt"

type object struct {
	class *Class
}

type RObject interface {
	Class() *Class
}

var _ RObject = &Object{}

type Object struct {
	object
}

func (obj *object) Class() *Class {
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
