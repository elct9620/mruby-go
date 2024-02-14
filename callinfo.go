package mruby

type callinfo struct {
	numArgs     int
	methodId    Symbol
	stackOffset int
	targetClass RClass
	proc        RProc
}

func (ci *callinfo) TargetClass() RClass {
	return ci.targetClass
}
