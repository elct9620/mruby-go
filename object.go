package mruby

import "fmt"

type RBasic struct {
}

type RObject struct {
	RBasic
}

func initObject(mrb *State) {
	mrb.falseClass = NewClass()
	mrb.trueClass = NewClass()

	mrb.objectClass.DefineMethod("puts", &Method{
		Function: func(mrb *State, recv Value) Value {
			args := mrb.GetArgv()
			fmt.Println(args...)
			return args[0]
		},
	})
}
