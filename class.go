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

func (mrb *State) FindMethod(recv Value, mid string) *Method {
	if m, ok := methods[mid]; ok {
		return m
	}

	return nil
}
