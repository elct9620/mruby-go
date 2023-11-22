package mruby

type context struct {
	ciCursor  int
	callinfos []*callinfo
}

func (ctx *context) GetCallinfo() *callinfo {
	return ctx.callinfos[len(ctx.callinfos)-1]
}
