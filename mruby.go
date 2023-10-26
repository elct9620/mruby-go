package mruby

type (
	Value = any
)

type Function func(*State, Value) Value

type Method struct {
	Function
}

type State struct {
}

func New() *State {
	return &State{}
}
