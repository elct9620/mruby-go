package mruby

import (
	"fmt"
	"encoding/binary"
	"errors"
)

var (
	ErrIRepNotFound = errors.New("irep not found")
)

func (mrb *State) TopRun(proc RProc, self Value) (Value, error) {
	mrb.context.callinfo.Push(&callinfo{})

	return mrb.VmRun(proc, self)
}

func (mrb *State) VmRun(proc RProc, self Value) (Value, error) {
	if mrb.context.stack == nil {
		mrb.context.stack = make([]Value, StackInitSize)
		mrb.context.stackBase = 0
		mrb.context.stackEnd = StackInitSize - 1
	}

	mrb.context.stack[0] = mrb.topSelf

	return mrb.VmExec(proc)
}

func (mrb *State) VmExec(proc RProc) (Value, error) {
	ir, ok := proc.Body().(*iRep)
	if !ok {
		return nil, ErrIRepNotFound
	}

	ci := mrb.context.GetCallinfo()
	offset := ci.stackOffset
	regs := mrb.context.stack

	for {
		opCode := ir.iSeq.Next()

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

			method := mrb.VmFindMethod(recv, ci.targetClass, ci.methodId)

			if method == nil {
				regs[offset+int(a)] = nil
				mrb.PopCallinfo()
				break
			}

			regs[offset+int(a)] = method.Call(mrb, recv)
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

func (state *State) PushCallinfo(mid Symbol, pushStack int, argc byte, targetClass *Class) *callinfo {
	ctx := state.context
	prevCi := ctx.GetCallinfo()

	callinfo := &callinfo{
		methodId:    mid,
		stackOffset: prevCi.stackOffset + pushStack,
		numArgs:     int(argc & 0xf),
		targetClass: targetClass,
	}
	ctx.callinfo.Push(callinfo)

	return callinfo
}

func (state *State) PopCallinfo() {
	ctx := state.context
	ctx.callinfo.Pop()
}
