package mruby

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/elct9620/mruby-go/insn"
	"github.com/elct9620/mruby-go/op"
)

var (
	ErrIRepNotFound = errors.New("irep not found")
)

func (mrb *State) TopRun(proc RProc, self Value) (Value, error) {
	mrb.context.callinfo.Push(&callinfo{})

	return mrb.VmRun(proc, self)
}

func (mrb *State) VmRun(proc RProc, self Value) (Value, error) {
	irep, ok := proc.Body().(*insn.Representation)
	if !ok {
		return nil, ErrIRepNotFound
	}

	if mrb.context.stack == nil {
		mrb.context.stack = make([]Value, StackInitSize)
	}

	mrb.context.stack[0] = mrb.topSelf

	return mrb.VmExec(proc, irep.Sequence().Clone())
}

func (mrb *State) VmExec(proc RProc, code *insn.Sequence) (Value, error) {
	rep, ok := proc.Body().(*insn.Representation)
	if !ok {
		return nil, ErrIRepNotFound
	}

	ctx := mrb.context

	for {
		opCode := code.Operation()

		switch opCode {
		case op.Move:
			ctx.Set(int(code.ReadB()), ctx.Get(int(code.ReadB())))
		case op.LoadI:
			a := code.ReadB()
			b := code.ReadB()
			ctx.Set(int(a), int(b))
		case op.LoadINeg:
			a := code.ReadB()
			b := code.ReadB()
			ctx.Set(int(a), -int(b))
		case op.LoadI__1, op.LoadI_0, op.LoadI_1, op.LoadI_2, op.LoadI_3, op.LoadI_4, op.LoadI_5, op.LoadI_6, op.LoadI_7:
			a := code.ReadB()
			ctx.Set(int(a), int(opCode)-int(op.LoadI_0))
		case op.LoadT, op.LoadF:
			a := code.ReadB()
			ctx.Set(int(a), opCode == op.LoadT)
		case op.LoadI16:
			a := code.ReadB()
			b := code.ReadS()
			ctx.Set(int(a), int(int16(binary.BigEndian.Uint16(b))))
		case op.LoadI32:
			a := code.ReadB()
			b := code.ReadW()
			ctx.Set(int(a), int(int32(binary.BigEndian.Uint32(b))))
		case op.LoadSym:
			a := code.ReadB()
			b := code.ReadB()
			ctx.Set(int(a), rep.Symbol(b))
		case op.LoadNil:
			a := code.ReadB()
			ctx.Set(int(a), nil)
		case op.GetConst:
			a := code.ReadB()
			b := code.ReadB()
			ctx.Set(int(a), mrb.VmGetConst(rep.Symbol(b)))
		case op.SetConst:
			a := code.ReadB()
			b := code.ReadB()
			mrb.VmSetConst(rep.Symbol(b), ctx.Get(int(a)))
		case op.SelfSend, op.Send, op.SendB:
			a := code.ReadB()
			b := code.ReadB()
			c := code.ReadB()

			if opCode == op.SelfSend {
				ctx.Set(int(a), ctx.Get(0))
				opCode = op.Send //nolint:ineffassign
			}

			mid := rep.Symbol(b)

			ci := mrb.callinfoPush(mid, int(a), c, nil)
			recv := ctx.Get(0)
			ci.targetClass = mrb.Class(recv)

			method := mrb.VmFindMethod(recv, ci.targetClass, ci.methodId)

			if method == nil {
				ctx.Set(int(a), nil)
				mrb.callinfoPop()
				break
			}

			ctx.Set(0, method.Call(mrb, recv))
			mrb.callinfoPop()
		case op.String:
			a := code.ReadB()
			b := code.ReadB()

			ctx.Set(int(a), rep.PoolValue(b))
		case op.Return:
			a := code.ReadB()
			return ctx.Get(int(a)), nil
		case op.StrCat:
			a := code.ReadB()
			ctx.Set(int(a), fmt.Sprintf("%v%v", ctx.Get(int(a)), ctx.Get(int(a)+1)))
		case op.Class:
			a := code.ReadB()
			b := code.ReadB()

			base := ctx.Get(int(a))
			super := ctx.Get(int(a) + 1)
			id := rep.Symbol(b)

			if base == nil {
				base = mrb.ObjectClass
			}

			class, err := mrb.vmDefineClass(base, super, id)
			if err != nil {
				return nil, err
			}

			ctx.Set(int(a), NewObjectValue(class))
		default:
			return nil, fmt.Errorf("opcode %d not implemented", opCode)
		}
	}
}

func (state *State) callinfoPush(mid Symbol, pushStack int, argc byte, targetClass *Class) *callinfo {
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

func (state *State) callinfoPop() {
	ctx := state.context
	ctx.callinfo.Pop()
}
