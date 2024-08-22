package mruby

import "fmt"

type RBasic interface {
	Class() RClass
	Flags() uint32
}

type RObject interface {
	RBasic
	ivPut(Symbol, Value)
	ivGet(Symbol) Value
}

var _ RBasic = &Object{}
var _ RObject = &Object{}

type Object struct {
	class RClass
	flags uint32
	iv
}

func (obj *Object) Class() RClass {
	return obj.class
}

func (obj *Object) Flags() uint32 {
	return obj.flags
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

func objectPuts(mrb *State, recv Value) Value {
	args := mrb.GetArgv()
	fmt.Println(args...)
	return args[0]
}

func objectRaise(mrb *State, recv Value) Value {
	args := mrb.GetArgv()
	mrb.Raise(nil, fmt.Sprint(args...))

	return nil
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

	mrb.DefineMethodId(mrb.ObjectClass, _raise(mrb), objectRaise)
	return nil
}
