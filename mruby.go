package mruby

import (
	"github.com/elct9620/mruby-go/stack"
)

type (
	Value = any
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
	targetClass *Class
}

type State struct {
	context *context

	ObjectClass  *Class
	ModuleClass  *Class
	ClassClass   *Class
	FalseClass   *Class
	TrueClass    *Class
	KernelModule *Class

	topSelf *RObject

	symbolTable map[string]Symbol
	symbolIndex int
}

func New() *State {
	state := &State{
		context: &context{
			callinfo: stack.New[*callinfo](),
		},
		topSelf:     &RObject{},
		symbolTable: make(map[string]Symbol),
		symbolIndex: 0,
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
