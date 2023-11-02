package mruby

import (
	"encoding/binary"
	"fmt"
)

const nullSymbolLength = 0xffff

var _ executable = &iRep{}

type iRepReaderFn func(*iRep, *Reader) error

var iRepReaders = []iRepReaderFn{
	readIRepHeader,
	readISeq,
	readPoolValues,
	readSyms,
}

type iRep struct {
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

func newIRep(r *Reader) (*iRep, error) {
	iRep := &iRep{}

	for _, loader := range iRepReaders {
		err := loader(iRep, r)
		if err != nil {
			return nil, err
		}
	}

	return iRep, nil
}

func (ir *iRep) Execute(state *State) (Value, error) {
	regs := make([]Value, ir.nRegs)

	for {
		opCode := ir.iSeq.ReadCode()

		switch opCode {
		case opMove:
			regs[ir.iSeq.ReadB()] = regs[ir.iSeq.ReadB()]
		case opLoadI:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			regs[a] = int(b)
		case opLoadINeg:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			regs[a] = -int(b)
		case opLoadI__1, opLoadI_0, opLoadI_1, opLoadI_2, opLoadI_3, opLoadI_4, opLoadI_5, opLoadI_6, opLoadI_7:
			a := ir.iSeq.ReadB()
			regs[a] = int(opCode) - int(opLoadI_0)
		case opLoadT, opLoadF:
			a := ir.iSeq.ReadB()
			regs[a] = opCode == opLoadT
		case opLoadI16:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadS()
			regs[a] = int(int16(binary.BigEndian.Uint16(b)))
		case opLoadI32:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadW()
			regs[a] = int(int32(binary.BigEndian.Uint32(b)))
		case opLoadSym:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			regs[a] = ir.syms[b]
		case opSelfSend:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			c := ir.iSeq.ReadB()

			ci := &Callinfo{
				NumArgs:  int(c & 0xf),
				MethodId: ir.syms[b],
				Stack:    []Value{nil},
			}

			ci.Stack = append(ci.Stack, regs[int(a)+1:int(a)+ci.NumArgs+1]...)
			state.Context.Callinfo = ci

			recv := regs[0]
			method := findMethod(state, recv, ci.MethodId)

			if method == nil {
				regs[a] = nil
				break
			}

			regs[a] = method.Function(state, recv)
		case opString:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()

			regs[a] = ir.poolValue[b]
		case opReturn:
			a := ir.iSeq.ReadB()
			return regs[a], nil
		case opStrCat:
			a := ir.iSeq.ReadB()
			regs[a] = fmt.Sprintf("%v%v", regs[a], regs[a+1])
		default:
			return nil, fmt.Errorf("opcode %d not implemented", opCode)
		}
	}
}

func readIRepHeader(ir *iRep, r *Reader) (err error) {
	fields := []any{
		&ir.nLocals,
		&ir.nRegs,
		&ir.rLen,
		&ir.cLen,
		&ir.iLen,
	}

	for _, field := range fields {
		err = r.ReadAs(field)
		if err != nil {
			return err
		}
	}

	return nil
}

func readISeq(ir *iRep, r *Reader) error {
	binary := make([]byte, ir.iLen)
	err := r.ReadAs(binary)
	if err != nil {
		return err
	}

	ir.iSeq = newInstructionSequence(binary)
	return nil
}

func readSyms(ir *iRep, r *Reader) error {
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
