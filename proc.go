package mruby

import "github.com/elct9620/mruby-go/insn"

type RProc interface {
	Body() any
}

var _ RProc = &proc{}
var _ RBasic = &proc{}

type proc struct {
	object
	body any
}

func (p *proc) Body() any {
	return p.body
}

func newMethodFromFunc(function Function) Method {
	return &method{
		Function: function,
	}
}

func newMethodFromProc(proc RProc) Method {
	return &method{
		RProc: proc,
	}
}

func (mrb *State) procNew(irep *insn.Representation) RProc {
	return &proc{
		body: irep,
	}
}
