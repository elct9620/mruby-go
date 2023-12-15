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

	ObjectClass  *RClass
	ModuleClass  *RClass
	ClassClass   *RClass
	FalseClass   *RClass
	TrueClass    *RClass
	KernelModule *RClass

	topSelf *RObject
}

func New() *State {
	state := &State{
		context: &context{
			callinfo: stack.New[*callinfo](),
		},
		topSelf: &RObject{},
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
