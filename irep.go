package mruby

import (
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
	cursor  int
}

func newIrep(r io.Reader) (*irep, error) {
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
	var a uint8
	regs := make([]value, ir.nRegs)

	for {
		opCode := ir.iSeq[ir.cursor]
		ir.cursor++

		switch opCode {
		case opLOADI__1, opLOADI_0, opLOADI_1, opLOADI_2, opLOADI_3, opLOADI_4, opLOADI_5, opLOADI_6, opLOADI_7:
			a = ir.iSeq[ir.cursor]
			regs[a] = int(opCode) - int(opLOADI_0)
			ir.cursor++
		case opRETURN:
			a = ir.iSeq[ir.cursor]
			return regs[a], nil
		default:
			return nil, fmt.Errorf("opcode %d not implemented", opCode)
		}
	}
}

func irepReadHeader(r io.Reader, ir *irep) (err error) {
	err = binaryRead(r, &ir.nLocals)
	if err != nil {
		return err
	}

	err = binaryRead(r, &ir.nRegs)
	if err != nil {
		return err
	}
	err = binaryRead(r, &ir.rLen)
	if err != nil {
		return err
	}
	err = binaryRead(r, &ir.cLen)
	if err != nil {
		return err
	}

	return binaryRead(r, &ir.iLen)
}

func irepReadISeq(r io.Reader, ir *irep) error {
	ir.iSeq = make([]code, ir.iLen)

	return binaryRead(r, ir.iSeq)
}
