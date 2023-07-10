package mruby

type Executable interface {
	Execute() Value
}

type Proc struct {
	Executable
}
