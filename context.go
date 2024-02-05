package mruby

import "github.com/elct9620/mruby-go/stack"

type context struct {
	stack    *stack.Stack[Value]
	callinfo *stack.Stack[*callinfo]
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
