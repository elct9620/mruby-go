package mruby

import "fmt"

type object struct {
	class *Class
}

type Basic interface {
	Class() *Class
}

var _ Basic = &RObject{}

type RObject struct {
	object
}

func (obj *object) Class() *Class {
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
