package mruby

import (
	"bytes"

	"github.com/elct9620/mruby-go/rite"
)

type Mrb struct {
	rite *rite.RITE
}

func NewFromString(code string) (*Mrb, error) {
	compiled, err := Compile(bytes.NewBufferString(code))
	if err != nil {
		return nil, err
	}

	rite, err := rite.Load(bytes.NewBuffer(compiled))
	if err != nil {
		return nil, err
	}

	return &Mrb{
		rite: rite,
	}, nil
}

func (mrb *Mrb) Header() rite.Header {
	return mrb.rite.Header()
}
