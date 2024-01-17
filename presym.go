package mruby

type predefineSymbolFunc func(mrb *State) Symbol

func predefineSymbol(name string) predefineSymbolFunc {
	return func(mrb *State) Symbol {
		return mrb.Intern(name)
	}
}

var (
	_classname   = predefineSymbol("__classname__")
	_BasicObject = predefineSymbol("BasicObject")
	_Object      = predefineSymbol("Object")
	_Module      = predefineSymbol("Module")
	_Class       = predefineSymbol("Class")
)
