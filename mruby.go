package mruby

type (
	Value  = any
	Symbol = string
)

type Function func(*State, Value) Value

type Method struct {
	Function
}

type Callinfo struct {
	NumArgs  int
	MethodId Symbol
	Stack    []Value
}

type Context struct {
	Callinfo *Callinfo
}

type State struct {
	Context *Context
}

func New() *State {
	return &State{
		Context: &Context{},
	}
}

func (s *State) GetArgc() int {
	return s.Context.Callinfo.NumArgs
}

func (s *State) GetArgv() []Value {
	return s.Context.Callinfo.Stack[1:]
}
