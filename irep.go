package mruby

import (
	"encoding/binary"
	"fmt"
	"io"
)

var _ executable = &irep{}

type irep struct {
	nLocals   uint16
	nRegs     uint16
	rLen      uint16
	cLen      uint16
	iLen      uint32
	pLen      uint16
	iSeq      []code
	poolValue []Value
	cursor    int
}

func newIrep(r io.Reader) (*irep, error) {
	irep := &irep{}

	err := irepReadHeader(irep, r)
	if err != nil {
		return nil, err
	}

	err = irepReadISeq(irep, r)
	if err != nil {
		return nil, err
	}

	err = readPoolValues(irep, r)
	if err != nil {
		return nil, err
	}

	return irep, nil
}

func (ir *irep) Execute(state *State) (Value, error) {
	regs := make([]Value, ir.nRegs)

	for {
		opCode := ir.iSeq[ir.cursor]
		ir.cursor++

		switch opCode {
		case opLOADI:
			a := ir.readB()
			b := ir.readB()
			regs[a] = int(b)
		case opLOADINEG:
			a := ir.readB()
			b := ir.readB()
			regs[a] = -int(b)
		case opLOADI__1, opLOADI_0, opLOADI_1, opLOADI_2, opLOADI_3, opLOADI_4, opLOADI_5, opLOADI_6, opLOADI_7:
			a := ir.readB()
			regs[a] = int(opCode) - int(opLOADI_0)
		case opLOADT, opLOADF:
			a := ir.readB()
			regs[a] = opCode == opLOADT
		case opLOADI16:
			a := ir.readB()
			b := ir.readS()
			regs[a] = int(int16(binary.BigEndian.Uint16(b)))
		case opLOADI32:
			a := ir.readB()
			b := ir.readW()
			regs[a] = int(int32(binary.BigEndian.Uint32(b)))
		case opSTRING:
			a := ir.readB()
			b := ir.readB()

			regs[a] = ir.poolValue[b]
		case opRETURN:
			a := ir.readB()
			return regs[a], nil
		default:
			return nil, fmt.Errorf("opcode %d not implemented", opCode)
		}
	}
}

func (ir *irep) readB() byte {
	b := ir.iSeq[ir.cursor]
	ir.cursor++
	return b
}

func (ir *irep) readS() []byte {
	s := ir.iSeq[ir.cursor : ir.cursor+2]
	ir.cursor += 2
	return s
}

func (ir *irep) readW() []byte {
	w := ir.iSeq[ir.cursor : ir.cursor+4]
	ir.cursor += 4
	return w
}

func irepReadHeader(ir *irep, r io.Reader) (err error) {
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

func irepReadISeq(ir *irep, r io.Reader) error {
	ir.iSeq = make([]code, ir.iLen)

	return binaryRead(r, ir.iSeq)
}
