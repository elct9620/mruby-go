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

func anyToString(mrb *State, obj Value) Value {
	className := mrb.ObjectClassName(obj)
	return fmt.Sprintf("#<%s:%p>", className, obj)
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

	return nil
}
