package mruby

import "fmt"

type object struct {
	class RClass
}

type RObject interface {
	Class() RClass
	ivPut(Symbol, Value)
	ivGet(Symbol) Value
}

var _ RObject = &Object{}

type Object struct {
	object
	iv
}

func (obj *object) Class() RClass {
	return obj.class
}

func (obj *Object) ivPut(sym Symbol, val Value) {
	if obj.iv == nil {
		obj.iv = make(iv)
	}

	obj.iv.Put(sym, val)
}

func (obj *Object) ivGet(sym Symbol) Value {
	if obj.iv == nil {
		return nil
	}

	return obj.iv.Get(sym)
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

			return &Object{object{super}, nil}
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
