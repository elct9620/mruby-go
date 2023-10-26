package mruby

import "fmt"

var methods = map[string]*Method{
	"puts": {
		Function: func(mrb *State, recv Value) Value {
			args, ok := recv.([]Value)
			if !ok {
				return nil
			}

			fmt.Println(args...)
			return args[0]
		},
	},
}

func findMethod(mrb *State, recv Value, mid string) *Method {
	if m, ok := methods[mid]; ok {
		return m
	}

	return nil
}
