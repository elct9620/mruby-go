package mruby

import (
	"encoding/binary"
	"fmt"
)

const nullSymbolLength = 0xffff

var _ executable = &iRep{}

type iRepReaderFn func(*State, *iRep, *Reader) error

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
	syms      []Symbol
}

func newIRep(mrb *State, r *Reader) (*iRep, error) {
	iRep := &iRep{}

	for _, loader := range iRepReaders {
		err := loader(mrb, iRep, r)
		if err != nil {
			return nil, err
		}
	}

	return iRep, nil
}

func (ir *iRep) Execute(mrb *State) (Value, error) {
	ci := mrb.context.GetCallinfo()
	offset := ci.stackOffset
	regs := mrb.context.stack

	for {
		opCode := ir.iSeq.ReadCode()

		switch opCode {
		case opMove:
			regs[offset+int(ir.iSeq.ReadB())] = regs[offset+int(ir.iSeq.ReadB())]
		case opLoadI:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			regs[offset+int(a)] = int(b)
		case opLoadINeg:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			regs[offset+int(a)] = -int(b)
		case opLoadI__1, opLoadI_0, opLoadI_1, opLoadI_2, opLoadI_3, opLoadI_4, opLoadI_5, opLoadI_6, opLoadI_7:
			a := ir.iSeq.ReadB()
			regs[offset+int(a)] = int(opCode) - int(opLoadI_0)
		case opLoadT, opLoadF:
			a := ir.iSeq.ReadB()
			regs[offset+int(a)] = opCode == opLoadT
		case opLoadI16:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadS()
			regs[offset+int(a)] = int(int16(binary.BigEndian.Uint16(b)))
		case opLoadI32:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadW()
			regs[offset+int(a)] = int(int32(binary.BigEndian.Uint32(b)))
		case opLoadSym:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			regs[offset+int(a)] = ir.syms[b]
		case opLoadNil:
			a := ir.iSeq.ReadB()
			regs[offset+int(a)] = nil
		case opGetConst:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			regs[a] = mrb.VmGetConst(ir.syms[b])
		case opSetConst:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			mrb.VmSetConst(ir.syms[b], regs[a])
		case opSelfSend, opSend, opSendB:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()
			c := ir.iSeq.ReadB()

			if opCode == opSelfSend {
				regs[offset+int(a)] = regs[offset]
				opCode = opSend //nolint:ineffassign
			}

			mid := ir.syms[b]

			ci := mrb.PushCallinfo(mid, int(a), c, nil)
			ci.stack = append(ci.stack, regs[int(a)+1:int(a)+ci.numArgs+1]...)

			recv := regs[offset]
			ci.targetClass = mrb.ClassOf(recv)
			method := mrb.FindMethod(recv, ci.targetClass, ci.methodId)

			if method == nil {
				regs[offset+int(a)] = nil
				mrb.PopCallinfo()
				break
			}

			regs[offset+int(a)] = method.Function(mrb, recv)
			mrb.PopCallinfo()
		case opString:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()

			regs[offset+int(a)] = ir.poolValue[b]
		case opReturn:
			a := ir.iSeq.ReadB()
			return regs[offset+int(a)], nil
		case opStrCat:
			a := ir.iSeq.ReadB()
			regs[offset+int(a)] = fmt.Sprintf("%v%v", regs[offset+int(a)], regs[offset+int(a)+1])
		case opClass:
			a := ir.iSeq.ReadB()
			b := ir.iSeq.ReadB()

			base := regs[offset+int(a)]
			super := regs[offset+int(a)+1]
			id := ir.syms[b]

			if base == nil {
				base = mrb.ObjectClass
			}

			class, err := mrb.vmDefineClass(base, super, id)
			if err != nil {
				return nil, err
			}

			regs[offset+int(a)] = NewObjectValue(class)
		default:
			return nil, fmt.Errorf("opcode %d not implemented", opCode)
		}
	}
}

func readIRepHeader(mrb *State, ir *iRep, r *Reader) (err error) {
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

func readISeq(mrb *State, ir *iRep, r *Reader) error {
	binary := make([]byte, ir.iLen)
	err := r.ReadAs(binary)
	if err != nil {
		return err
	}

	ir.iSeq = newInstructionSequence(binary)
	return nil
}

func readSyms(mrb *State, ir *iRep, r *Reader) error {
	err := r.ReadAs(&ir.sLen)
	if err != nil {
		return err
	}

	ir.syms = make([]Symbol, ir.sLen)

	for i := uint16(0); i < ir.sLen; i++ {
		strLen, err := r.ReadUint16()
		if err != nil {
			return err
		}

		if strLen == nullSymbolLength {
			ir.syms[i] = 0
			continue
		}

		str, err := r.ReadString(int(strLen) + 1)
		if err != nil {
			return err
		}

		ir.syms[i] = mrb.Intern(string(str))
	}

	return nil
}
