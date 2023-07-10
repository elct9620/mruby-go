package mruby

type Executable interface {
	Execute(state *State) (value, error)
}

type Proc struct {
	Executable
}
