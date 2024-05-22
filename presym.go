package mruby

type predefineSymbolFunc func(mrb *State) Symbol

func predefineSymbol(name string) predefineSymbolFunc {
	return func(mrb *State) Symbol {
		return mrb.Intern(name)
	}
}

var (
	_classname   = predefineSymbol("__classname__")
	_attached    = predefineSymbol("__attached__")
	_new         = predefineSymbol("new")
	_BasicObject = predefineSymbol("BasicObject")
	_Object      = predefineSymbol("Object")
	_Module      = predefineSymbol("Module")
	_Class       = predefineSymbol("Class")
	_Kernel      = predefineSymbol("Kernel")
	_TrueClass   = predefineSymbol("TrueClass")
	_FalseClass  = predefineSymbol("FalseClass")
	_ArrayClass  = predefineSymbol("Array")
	_join        = predefineSymbol("join")
	_raise       = predefineSymbol("raise")
)
