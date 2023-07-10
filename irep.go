package mruby

import (
	"encoding/binary"
	"fmt"
	"io"
)

var _ Executable = &IREP{}

type IREP struct {
	NLocals uint16
	NRegs   uint16
	RLen    uint16
	CLen    uint16
	ILen    uint32
}

func newIREP(r io.Reader) (*IREP, error) {
	irep := &IREP{}

	err := binary.Read(r, riteByteOrder, irep)
	if err != nil {
		return nil, err
	}

	return irep, nil
}

func (ir *IREP) Execute(state *State) (Value, error) {
	return fmt.Sprintf(
		"nlocals = %d, nregs = %d, rlen = %d, clen = %d, ilen = %d",
		ir.NLocals,
		ir.NRegs,
		ir.RLen,
		ir.CLen,
		ir.ILen,
	), nil
}
