package mruby

type callinfo struct {
	numArgs     int
	methodId    Symbol
	stackOffset int
	targetClass RClass
	proc        RProc
	pc          int // insn.Sequence cursor
}

func (ci *callinfo) TargetClass() RClass {
	return ci.targetClass
}

func (ci *callinfo) SetSequnceCursor(pc int) {
	ci.pc = pc
}

func (ci *callinfo) GetSequnceCursor() int {
	return ci.pc
}

func (ci *callinfo) Proc() RProc {
	return ci.proc
}
