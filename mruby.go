package mruby

import (
	"github.com/elct9620/mruby-go/stack"
)

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
	stackOffset int
	stack       []Value
	targetClass *RClass
}

type State struct {
	context *context

	falseClass  *RClass
	trueClass   *RClass
	objectClass *RClass

	topSelf *RObject
}

func New() *State {
	state := &State{
		context: &context{
			callinfo: stack.New[*callinfo](),
		},
		objectClass: NewClass(),
		topSelf:     &RObject{},
	}

	initCore(state)

	return state
}

func (s *State) GetArgc() int {
	return s.context.GetCallinfo().numArgs
}

func (s *State) GetArgv() []Value {
	ci := s.context.GetCallinfo()
	begin := ci.stackOffset + 1
	end := begin + ci.numArgs

	return s.context.stack[begin:end]
}

func initCore(mrb *State) {
	initObject(mrb)
}
