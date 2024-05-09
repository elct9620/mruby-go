package mruby

type ivTable map[Symbol]Value
type iv = ivTable

func (iv ivTable) Get(sym Symbol) Value {
	return iv[sym]
}

func (iv ivTable) Put(sym Symbol, val Value) {
	iv[sym] = val
}

func (mrb *State) ObjectInstanceVariableSetForce(obj RObject, name Symbol, val Value) {
	obj.ivPut(name, val)
}

func (mrb *State) ObjectInstanceVariableGet(obj RObject, name Symbol) Value {
	return obj.ivGet(name)
}

func (mrb *State) ObjectInstanceVariableDefined(obj RObject, name Symbol) bool {
	return obj.ivGet(name) != nil
}

func (mrb *State) DefineConstById(klass RClass, name Symbol, val Value) {
	mrb.ObjectInstanceVariableSetForce(klass, name, val)
}

func (mrb *State) VmGetConst(sym Symbol) Value {
	klass := mrb.ObjectClass
	return mrb.ObjectInstanceVariableGet(klass, sym)
}

func (mrb *State) VmSetConst(sym Symbol, val Value) {
	// NOTE: Find class in current context
	klass := mrb.ObjectClass
	mrb.ObjectInstanceVariableSetForce(klass, sym, val)
}
