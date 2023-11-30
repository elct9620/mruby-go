package mruby

import "github.com/elct9620/mruby-go/stack"

type (
	Value  = any
	Symbol = string
)

type Function func(*State, Value) Value

type Method struct {
	Function
}

type callinfo struct {
	numArgs     int
	methodId    Symbol
	stack       []Value
	targetClass *RClass
}

type State struct {
	context *context

	falseClass *RClass

	topSelf *RObject
}

func New() *State {
	return &State{
		context: &context{
			callinfo: stack.New[*callinfo](),
		},
		falseClass: &RClass{},
		topSelf:    &RObject{},
	}
}

func (s *State) GetArgc() int {
	return s.context.GetCallinfo().numArgs
}

func (s *State) GetArgv() []Value {
	return s.context.GetCallinfo().stack[1:]
}
