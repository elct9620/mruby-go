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
	numArgs     int
	methodId    Symbol
	stack       []Value
	targetClass *RClass
}

type context struct {
	ciCursor  int
	callinfos []*callinfo
}

func (ctx *context) GetCallinfo() *callinfo {
	return ctx.callinfos[len(ctx.callinfos)-1]
}

type State struct {
	context *context

	falseClass *RClass
}

func New() *State {
	return &State{
		context: &context{
			ciCursor:  -1,
			callinfos: []*callinfo{},
		},
		falseClass: &RClass{},
	}
}

func (s *State) GetArgc() int {
	return s.context.GetCallinfo().numArgs
}

func (s *State) GetArgv() []Value {
	return s.context.GetCallinfo().stack[1:]
}
