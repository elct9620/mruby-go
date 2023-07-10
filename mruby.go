package mruby

type value = any

type state struct {
}

func New() *state {
	return &state{}
}
