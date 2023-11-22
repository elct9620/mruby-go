package mruby

import "github.com/elct9620/mruby-go/stack"

type context struct {
	callinfo *stack.Stack[*callinfo]
}

func (ctx *context) GetCallinfo() *callinfo {
	return ctx.callinfo.Peek()
}
