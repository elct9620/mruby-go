package mruby

type Value = any

type State struct {
}

func New() *State {
	return &State{}
}
