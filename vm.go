package mruby

func (state *State) PushCallinfo(mid string, argc byte, targetClass *RClass) *callinfo {
	ctx := state.context

	callinfo := &callinfo{
		methodId:    mid,
		numArgs:     int(argc & 0xf),
		stack:       []Value{nil},
		targetClass: targetClass,
	}

	ctx.ciCursor++

	if ctx.ciCursor >= len(state.context.callinfos) {
		ctx.callinfos = append(state.context.callinfos, callinfo)
	} else {
		ctx.callinfos[ctx.ciCursor] = callinfo
	}

	return callinfo
}

func (state *State) PopCallinfo() {
	ctx := state.context
	ctx.callinfos[ctx.ciCursor] = nil
	ctx.ciCursor--
}
