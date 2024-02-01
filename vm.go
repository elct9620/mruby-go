package mruby

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrIRepNotFound = errors.New("irep not found")
)

func (mrb *State) TopRun(proc RProc, self Value) (Value, error) {
	mrb.context.callinfo.Push(&callinfo{})

	return mrb.VmRun(proc, self)
}

func (mrb *State) VmRun(proc RProc, self Value) (Value, error) {
	ir, ok := proc.Body().(*iRep)
	if !ok {
		return nil, ErrIRepNotFound
	}

	if mrb.context.stack == nil {
		mrb.context.stack = make([]Value, StackInitSize)
	}

	mrb.context.stack[0] = mrb.topSelf

	return mrb.VmExec(proc, ir.iSeq.Clone())
}

func (mrb *State) VmExec(proc RProc, code *Code) (Value, error) {
	ir, ok := proc.Body().(*iRep)
	if !ok {
		return nil, ErrIRepNotFound
	}

	ci := mrb.context.GetCallinfo()
	offset := ci.stackOffset
	regs := mrb.context.stack

	for {
		opCode := code.Next()

		switch opCode {
		case opMove:
			regs[offset+int(code.ReadB())] = regs[offset+int(code.ReadB())]
		case opLoadI:
			a := code.ReadB()
			b := code.ReadB()
			regs[offset+int(a)] = int(b)
		case opLoadINeg:
			a := code.ReadB()
			b := code.ReadB()
			regs[offset+int(a)] = -int(b)
		case opLoadI__1, opLoadI_0, opLoadI_1, opLoadI_2, opLoadI_3, opLoadI_4, opLoadI_5, opLoadI_6, opLoadI_7:
			a := code.ReadB()
			regs[offset+int(a)] = int(opCode) - int(opLoadI_0)
		case opLoadT, opLoadF:
			a := code.ReadB()
			regs[offset+int(a)] = opCode == opLoadT
		case opLoadI16:
			a := code.ReadB()
			b := code.ReadS()
			regs[offset+int(a)] = int(int16(binary.BigEndian.Uint16(b)))
		case opLoadI32:
			a := code.ReadB()
			b := code.ReadW()
			regs[offset+int(a)] = int(int32(binary.BigEndian.Uint32(b)))
		case opLoadSym:
			a := code.ReadB()
			b := code.ReadB()
			regs[offset+int(a)] = ir.syms[b]
		case opLoadNil:
			a := code.ReadB()
			regs[offset+int(a)] = nil
		case opGetConst:
			a := code.ReadB()
			b := code.ReadB()
			regs[a] = mrb.VmGetConst(ir.syms[b])
		case opSetConst:
			a := code.ReadB()
			b := code.ReadB()
			mrb.VmSetConst(ir.syms[b], regs[a])
		case opSelfSend, opSend, opSendB:
			a := code.ReadB()
			b := code.ReadB()
			c := code.ReadB()

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
			a := code.ReadB()
			b := code.ReadB()

			regs[offset+int(a)] = ir.poolValue[b]
		case opReturn:
			a := code.ReadB()
			return regs[offset+int(a)], nil
		case opStrCat:
			a := code.ReadB()
			regs[offset+int(a)] = fmt.Sprintf("%v%v", regs[offset+int(a)], regs[offset+int(a)+1])
		case opClass:
			a := code.ReadB()
			b := code.ReadB()

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
