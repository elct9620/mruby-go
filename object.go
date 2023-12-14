package mruby

import "fmt"

type RBasic struct {
}

type RObject struct {
	RBasic
}

func initObject(mrb *State) {
	mrb.falseClass = mrb.NewClass(mrb.objectClass)
	mrb.trueClass = mrb.NewClass(mrb.objectClass)

	mrb.objectClass.DefineMethod("puts", &Method{
		Function: func(mrb *State, recv Value) Value {
			args := mrb.GetArgv()
			fmt.Println(args...)
			return args[0]
		},
	})
}
