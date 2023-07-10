package mruby

type executable interface {
	Execute(state *state) (value, error)
}

type proc struct {
	executable
}
