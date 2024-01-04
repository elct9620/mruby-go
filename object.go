package mruby

import "fmt"

type object struct {
	class *RClass
}

type RBasic interface {
	Class() *RClass
}

var _ RBasic = &RObject{}

type RObject struct {
	object
}

func (obj *object) Class() *RClass {
	return nil
}

func initObject(mrb *State) {
	mrb.FalseClass = mrb.NewClass(mrb.ObjectClass)
	mrb.TrueClass = mrb.NewClass(mrb.ObjectClass)

	mrb.ObjectClass.DefineMethod(mrb, "puts", &Method{
		Function: func(mrb *State, recv Value) Value {
			args := mrb.GetArgv()
			fmt.Println(args...)
			return args[0]
		},
	})
}
