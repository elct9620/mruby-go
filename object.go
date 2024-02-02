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

func initObject(mrb *State) (err error) {
	mrb.FalseClass, err = mrb.DefineClassId(_FalseClass(mrb), mrb.ObjectClass)
	if err != nil {
		return err
	}
	mrb.TrueClass, err = mrb.DefineClassId(_TrueClass(mrb), mrb.ObjectClass)
	if err != nil {
		return err
	}

	putsSym := mrb.Intern("puts")
	mrb.DefineMethodId(mrb.ObjectClass, putsSym, objectPuts)

	return nil
}

func objectPuts(mrb *State, recv Value) Value {
	args := mrb.GetArgv()
	fmt.Println(args...)
	return args[0]
}
