package mruby

type RProc interface {
	Execute(state *State) (Value, error)
}

type executable interface {
	Execute(state *State) (Value, error)
}

var _ RProc = &proc{}

type proc struct {
	executable
}

func newMethodFromFunc(function Function) Method {
	return &goMethod{
		Function: function,
	}
}
