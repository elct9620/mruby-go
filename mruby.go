package mruby

type (
	Value = any
	code  = uint8
)

type State struct {
}

func New() *State {
	return &State{}
}
