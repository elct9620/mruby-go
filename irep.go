package mruby

import (
	"encoding/binary"
	"fmt"
	"io"
)

var _ executable = &irep{}

type irep struct {
	NLocals uint16
	NRegs   uint16
	RLen    uint16
	CLen    uint16
	ILen    uint32
}

func newIREP(r io.Reader) (*irep, error) {
	irep := &irep{}

	err := binary.Read(r, riteByteOrder, irep)
	if err != nil {
		return nil, err
	}

	return irep, nil
}

func (ir *irep) Execute(state *state) (value, error) {
	return fmt.Sprintf(
		"nlocals = %d, nregs = %d, rlen = %d, clen = %d, ilen = %d",
		ir.NLocals,
		ir.NRegs,
		ir.RLen,
		ir.CLen,
		ir.ILen,
	), nil
}
