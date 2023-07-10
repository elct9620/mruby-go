package mruby

import (
	"encoding/binary"
	"fmt"
	"io"
)

var _ executable = &irep{}

type irep struct {
	nLocals uint16
	nRegs   uint16
	rLen    uint16
	cLen    uint16
	iLen    uint32
	iSeq    []code
}

func newIREP(r io.Reader) (*irep, error) {
	irep := &irep{}

	err := irepReadHeader(r, irep)
	if err != nil {
		return nil, err
	}

	err = irepReadISeq(r, irep)
	if err != nil {
		return nil, err
	}

	return irep, nil
}

func (ir *irep) Execute(state *state) (value, error) {
	return fmt.Sprintf(
		"nlocals = %d, nregs = %d, rlen = %d, clen = %d, ilen = %d",
		ir.nLocals,
		ir.nRegs,
		ir.rLen,
		ir.cLen,
		ir.iLen,
	), nil
}

func irepReadHeader(r io.Reader, ir *irep) (err error) {
	err = binary.Read(r, riteByteOrder, &ir.nLocals)
	if err != nil {
		return err
	}

	err = binary.Read(r, riteByteOrder, &ir.nRegs)
	if err != nil {
		return err
	}
	err = binary.Read(r, riteByteOrder, &ir.rLen)
	if err != nil {
		return err
	}
	err = binary.Read(r, riteByteOrder, &ir.cLen)
	if err != nil {
		return err
	}

	return binary.Read(r, riteByteOrder, &ir.iLen)
}

func irepReadISeq(r io.Reader, ir *irep) error {
	ir.iSeq = make([]code, ir.iLen)

	return binary.Read(r, riteByteOrder, ir.iSeq)
}
