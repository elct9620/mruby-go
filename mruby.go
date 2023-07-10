package mruby

type Value = any

func New() *State {
	return &State{}
}
