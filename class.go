package mruby

import "fmt"

var methods = map[string]*Method{
	"puts": {
		Function: func(mrb *State, recv Value) Value {
			args := mrb.GetArgv()
			fmt.Println(args...)
			return args[0]
		},
	},
}

type RClass struct {
	super *RClass
}

func (mrb *State) ClassOf(v Value) *RClass {
	// NOTE: Null Pointer fallback to FalseClass
	return mrb.falseClass
}

func (mrb *State) FindMethod(recv Value, mid string) *Method {
	if m, ok := methods[mid]; ok {
		return m
	}

	return nil
}
