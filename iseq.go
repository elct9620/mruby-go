package mruby

import "github.com/elct9620/mruby-go/op"

// mrb_code aka iseq
type Code struct {
	binary []byte
	cursor int
}

func NewCode(bytes []byte) *Code {
	return &Code{binary: bytes, cursor: 0}
}

func (code *Code) Next() op.Code {
	return code.ReadB()
}

func (code *Code) ReadB() byte {
	b := code.binary[code.cursor]
	code.cursor++
	return b
}

func (code *Code) ReadS() []byte {
	s := code.binary[code.cursor : code.cursor+2]
	code.cursor += 2
	return s
}

func (code *Code) ReadW() []byte {
	w := code.binary[code.cursor : code.cursor+4]
	code.cursor += 4
	return w
}

func (code *Code) Clone() *Code {
	return &Code{binary: code.binary, cursor: code.cursor}
}
