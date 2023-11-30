package mruby

func (state *State) PushCallinfo(mid string, pushStack int, argc byte, targetClass *RClass) *callinfo {
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
