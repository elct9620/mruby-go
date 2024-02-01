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
