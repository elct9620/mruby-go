package mruby

type RProc interface {
	Body() any
}

var _ RProc = &proc{}

type proc struct {
	body any
}

func (p *proc) Body() any {
	return p.body
}

func newMethodFromFunc(function Function) Method {
	return &goMethod{
		Function: function,
	}
}
