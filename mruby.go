package mruby

import (
	"github.com/elct9620/mruby-go/stack"
)

const (
	CallinfoInitSize = 32
	StackInitSize    = 128
)

type (
	Value = any
)

type State struct {
	context *context

	ObjectClass  RClass
	ModuleClass  RClass
	ClassClass   RClass
	FalseClass   RClass
	TrueClass    RClass
	KernelModule RClass

	topSelf RObject

	symbolTable map[string]Symbol
	symbolIndex int
}

func New() (*State, error) {
	state := &State{
		context: &context{
			callinfo: stack.New[*callinfo](CallinfoInitSize),
		},
		symbolTable: make(map[string]Symbol),
		symbolIndex: 0,
	}

	err := initCore(state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (s *State) GetArgc() int {
	return s.context.GetCallinfo().numArgs
}

func (s *State) GetArgv() []Value {
	ci := s.context.GetCallinfo()

	return s.context.Slice(1, ci.numArgs)
}
