package mruby

type ivTable map[Symbol]Value
type iv = ivTable

func (iv ivTable) Get(sym Symbol) Value {
	return iv[sym]
}

func (iv ivTable) Set(sym Symbol, val Value) {
	iv[sym] = val
}

func (mrb *State) VmGetConst(sym Symbol) Value {
	klass := mrb.ObjectClass
	return klass.Get(sym)
}

func (mrb *State) VmSetConst(sym Symbol, val Value) {
	// NOTE: Find class in current context
	klass := mrb.ObjectClass
	klass.Set(sym, val)
}
