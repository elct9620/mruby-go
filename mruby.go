package mruby

type Value = any

type State struct {
}

func New() *State {
	return &State{}
}

func (s *State) Run(proc *Proc) Value {
	return proc.Execute()
}
