package mruby

import "github.com/elct9620/mruby-go/stack"

const StackInitSize = 128

type context struct {
	stack    []Value
	callinfo *stack.Stack[*callinfo]
}

func (ctx *context) GetCallinfo() *callinfo {
	return ctx.callinfo.Peek()
}

func (ctx *context) Get(idx int) Value {
	offset := ctx.GetCallinfo().stackOffset
	return ctx.stack[offset+idx]
}

func (ctx *context) Slice(start, end int) []Value {
	offset := ctx.GetCallinfo().stackOffset
	return ctx.stack[offset+start : offset+start+end]
}

func (ctx *context) Set(idx int, v Value) {
	offset := ctx.GetCallinfo().stackOffset
	ctx.stack[offset+idx] = v
}
