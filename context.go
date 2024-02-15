package mruby

type CallinfoStack interface {
	Cursor() int
	Push(*callinfo)
	Pop() *callinfo
	Peek() *callinfo
}

type ContextStack interface {
	Get(int) Value
	Slice(int, int) []Value
	Set(int, Value)
}

type context struct {
	stack    ContextStack
	callinfo CallinfoStack
}

func (ctx *context) IsTop() bool {
	return ctx.callinfo.Cursor() == 0
}

func (ctx *context) GetCallinfo() *callinfo {
	return ctx.callinfo.Peek()
}

func (ctx *context) Get(idx int) Value {
	offset := ctx.GetCallinfo().stackOffset
	return ctx.stack.Get(offset + idx)
}

func (ctx *context) Slice(start, end int) []Value {
	offset := ctx.GetCallinfo().stackOffset
	return ctx.stack.Slice(offset+start, offset+start+end)
}

func (ctx *context) Set(idx int, v Value) {
	offset := ctx.GetCallinfo().stackOffset
	ctx.stack.Set(offset+idx, v)
}
