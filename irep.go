package mruby

import (
	"encoding/binary"
	"fmt"
)

const nullSymbolLength = 0xffff

var _ executable = &irep{}

type irep struct {
	nLocals   uint16
	nRegs     uint16
	rLen      uint16
	cLen      uint16
	iLen      uint32
	pLen      uint16
	sLen      uint16
	iSeq      *instructionSequence
	poolValue []Value
	syms      []string
}

func newIrep(r *Reader) (*irep, error) {
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

	err = readSyms(irep, r)
	if err != nil {
		return nil, err
	}

	return irep, nil
}

func (ir *irep) Execute(state *State) (Value, error) {
	regs := make([]Value, ir.nRegs)

	for {
		opCode := ir.iSeq.ReadCode()

		switch opCode {
		case opLOADI:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			regs[a] = int(b)
		case opLOADINEG:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			regs[a] = -int(b)
		case opLOADI__1, opLOADI_0, opLOADI_1, opLOADI_2, opLOADI_3, opLOADI_4, opLOADI_5, opLOADI_6, opLOADI_7:
			a := ir.iSeq.ReadB()
			regs[a] = int(opCode) - int(opLOADI_0)
		case opLOADT, opLOADF:
			a := ir.iSeq.ReadB()
			regs[a] = opCode == opLOADT
		case opLOADI16:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadS()
			regs[a] = int(int16(binary.BigEndian.Uint16(b)))
		case opLOADI32:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadW()
			regs[a] = int(int32(binary.BigEndian.Uint32(b)))
		case opLOADSYM:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			regs[a] = ir.syms[b]
		case opSTRING:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()

			regs[a] = ir.poolValue[b]
		case opRETURN:
			a := ir.iSeq.ReadB()
			return regs[a], nil
		default:
			return nil, fmt.Errorf("opcode %d not implemented", opCode)
		}
	}
}

func irepReadHeader(ir *irep, r *Reader) (err error) {
	err = r.ReadAs(&ir.nLocals)
	if err != nil {
		return err
	}

	err = r.ReadAs(&ir.nRegs)
	if err != nil {
		return err
	}
	err = r.ReadAs(&ir.rLen)
	if err != nil {
		return err
	}
	err = r.ReadAs(&ir.cLen)
	if err != nil {
		return err
	}

	return r.ReadAs(&ir.iLen)
}

func irepReadISeq(ir *irep, r *Reader) error {
	binary := make([]byte, ir.iLen)
	err := r.ReadAs(binary)
	if err != nil {
		return err
	}

	ir.iSeq = newInstructionSequence(binary)
	return nil
}

func readSyms(ir *irep, r *Reader) error {
	err := r.ReadAs(&ir.sLen)
	if err != nil {
		return err
	}

	ir.syms = make([]string, ir.sLen)

	for i := uint16(0); i < ir.sLen; i++ {
		strLen, err := r.ReadUint16()
		if err != nil {
			return err
		}

		if strLen == nullSymbolLength {
			continue
		}

		str := make([]byte, strLen)
		err = r.ReadAs(str)
		if err != nil {
			return err
		}

		ir.syms[i] = string(str)
	}

	return nil
}
