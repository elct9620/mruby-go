package mruby

import (
	"github.com/elct9620/mruby-go/stack"
)

type (
	Value = any
)

type callinfo struct {
	numArgs     int
	methodId    Symbol
	stackOffset int
	stack       []Value
	targetClass *Class
}

type State struct {
	context *context

	ObjectClass  *Class
	ModuleClass  *Class
	ClassClass   *Class
	FalseClass   *Class
	TrueClass    *Class
	KernelModule *Module

	topSelf *Object

	symbolTable map[string]Symbol
	symbolIndex int
}

func New() (*State, error) {
	state := &State{
		context: &context{
			callinfo: stack.New[*callinfo](),
		},
		topSelf:     &Object{},
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
	begin := ci.stackOffset + 1
	end := begin + ci.numArgs

	return s.context.stack[begin:end]
}
