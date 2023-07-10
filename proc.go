package mruby

type Executable interface {
	Execute(state *State) (Value, error)
}

type Proc struct {
	Executable
}
