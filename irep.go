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

	err := irepReadHeader(r, irep)
	if err != nil {
		return nil, err
	}

	err = irepReadISeq(r, irep)
	if err != nil {
		return nil, err
	}

	err = irepReadPool(r, irep)
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

func irepReadPool(r io.Reader, ir *irep) error {
	var pLen uint16
	err := binaryRead(r, &pLen)
	if err != nil {
		return err
	}

	var pType uint8
	for i := 0; i < int(pLen); i++ {
		err = binaryRead(r, &pType)
		if err != nil {
			return err
		}

		switch pType {
		case poolTypeString:
			var sLen uint16
			err = binaryRead(r, &sLen)
			if err != nil {
				return err
			}

			s := make([]byte, sLen+1)
			err = binaryRead(r, s)
			if err != nil {
				return err
			}

			ir.poolValue = append(ir.poolValue, string(s[0:sLen]))
		case poolTypeInt32:
		case poolTypeStaticString:
		case poolTypeInt64:
		case poolTypeFloat:
		case poolTypeBigInt:
		}

		ir.pLen = uint16(i + 1)
	}

	return nil
}
