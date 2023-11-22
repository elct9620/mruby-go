package mruby

func (state *State) PushCallinfo(mid string, argc byte, targetClass *RClass) *callinfo {
	ctx := state.context

	callinfo := &callinfo{
		methodId:    mid,
		numArgs:     int(argc & 0xf),
		stack:       []Value{nil},
		targetClass: targetClass,
	}
	ctx.callinfo.Push(callinfo)

	return callinfo
}

func (state *State) PopCallinfo() {
	ctx := state.context
	ctx.callinfo.Pop()
}
