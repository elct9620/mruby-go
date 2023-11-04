package mruby

type (
	Value  = any
	Symbol = string
)

type Function func(*State, Value) Value

type Method struct {
	Function
}

type callinfo struct {
	numArgs  int
	methodId Symbol
	stack    []Value
}

type context struct {
	callinfo *callinfo
}

type State struct {
	context *context
}

func New() *State {
	return &State{
		context: &context{},
	}
}

func (s *State) GetArgc() int {
	return s.context.callinfo.numArgs
}

func (s *State) GetArgv() []Value {
	return s.context.callinfo.stack[1:]
}
