package mruby

func (mrb *State) TopRun(proc RProc, self Value) (Value, error) {
	mrb.context.callinfo.Push(&callinfo{})

	return mrb.VmRun(proc, self)
}

func (mrb *State) VmRun(proc RProc, self Value) (Value, error) {
	if mrb.context.stack == nil {
		mrb.context.stack = make([]Value, StackInitSize)
		mrb.context.stackBase = 0
		mrb.context.stackEnd = StackInitSize - 1
	}

	mrb.context.stack[0] = mrb.topSelf

	return proc.Execute(mrb)
}

func (state *State) PushCallinfo(mid Symbol, pushStack int, argc byte, targetClass *Class) *callinfo {
	ctx := state.context
	prevCi := ctx.GetCallinfo()

	callinfo := &callinfo{
		methodId:    mid,
		stackOffset: prevCi.stackOffset + pushStack,
		numArgs:     int(argc & 0xf),
		targetClass: targetClass,
	}
	ctx.callinfo.Push(callinfo)

	return callinfo
}

func (state *State) PopCallinfo() {
	ctx := state.context
	ctx.callinfo.Pop()
}
