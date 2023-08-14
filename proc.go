package mruby

type executable interface {
	Execute(state *State) (Value, error)
}

type proc struct {
	executable
}
