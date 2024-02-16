package mruby

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/elct9620/mruby-go/insn"
	"github.com/elct9620/mruby-go/op"
	"github.com/elct9620/mruby-go/stack"
)

var (
	ErrIRepNotFound               = errors.New("irep not found")
	ErrNotPrimitiveTypeComparison = errors.New("not primitive type comparison")
)

func (mrb *State) TopRun(proc RProc, self Value) (Value, error) {
	if !mrb.context.IsTop() {
		mrb.callinfoPush(0, 0, mrb.ObjectClass, proc, nil, 0, 0)
	}

	return mrb.VmRun(proc, self)
}

func (mrb *State) VmRun(proc RProc, self Value) (Value, error) {
	irep, ok := proc.Body().(*insn.Representation)
	if !ok {
		return nil, ErrIRepNotFound
	}

	if mrb.context.stack == nil {
		mrb.context.stack = stack.New[Value](StackInitSize)
	}

	mrb.context.Set(0, mrb.topSelf)

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
		case op.Nop:
			ctx.SetSequenceCursor(code.Cursor())
			continue
		case op.Move:
			a := code.ReadB()
			b := code.ReadB()

			ctx.SetSequenceCursor(code.Cursor())
			ctx.Set(int(a), ctx.Get(int(b)))
		case op.LoadL:
			a := code.ReadB()
			b := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), rep.PoolValue(b))
		case op.LoadI:
			a := code.ReadB()
			b := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), int(b))
		case op.LoadINeg:
			a := code.ReadB()
			b := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), -int(b))
		case op.LoadI__1, op.LoadI_0, op.LoadI_1, op.LoadI_2, op.LoadI_3, op.LoadI_4, op.LoadI_5, op.LoadI_6, op.LoadI_7:
			a := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), int(opCode)-int(op.LoadI_0))
		case op.LoadT, op.LoadF:
			a := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), opCode == op.LoadT)
		case op.LoadI16:
			a := code.ReadB()
			b := code.ReadS()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), int(int16(binary.BigEndian.Uint16(b))))
		case op.LoadI32:
			a := code.ReadB()
			b := code.ReadS()
			c := code.ReadS()
			ctx.SetSequenceCursor(code.Cursor())

			v := append(b, c...)
			ctx.Set(int(a), int(int32(binary.BigEndian.Uint32(v))))
		case op.LoadSym:
			a := code.ReadB()
			b := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), rep.Symbol(b))
		case op.LoadNil:
			a := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), nil)
		case op.GetConst:
			a := code.ReadB()
			b := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), mrb.VmGetConst(rep.Symbol(b)))
		case op.SetConst:
			a := code.ReadB()
			b := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			mrb.VmSetConst(rep.Symbol(b), ctx.Get(int(a)))
		case op.Jmp:
			a := code.ReadS()
			ctx.SetSequenceCursor(code.Cursor())

			offset := int(int16(binary.BigEndian.Uint16(a)))
			code.Seek(code.Cursor() + offset)
		case op.JmpNot:
			a := code.ReadB()
			b := code.ReadS()
			ctx.SetSequenceCursor(code.Cursor())

			val := ctx.Get(int(a))

			if !Test(val) {
				offset := int(int16(binary.BigEndian.Uint16(b)))
				code.Seek(code.Cursor() + offset)
			}
		case op.SelfSend, op.Send, op.SendB:
			a := code.ReadB()
			b := code.ReadB()
			c := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			if opCode == op.SelfSend {
				ctx.Set(int(a), ctx.Get(0))
				opCode = op.Send //nolint:ineffassign
			}

			mid := rep.Symbol(b)

			ci := mrb.callinfoPush(int(a), 0, nil, nil, nil, mid, uint16(c))
			recv := ctx.Get(0)
			ci.targetClass = mrb.Class(recv)

			method := mrb.VmFindMethod(recv, ci.targetClass, ci.methodId)

			if method == nil {
				ctx.Set(int(a), nil)
				mrb.callinfoPop()
				break
			}

			if method.IsProc() {
				proc := method.Proc()
				ci.proc = proc

				if !proc.IsGoFunction() {
					nirep := proc.Body().(*insn.Representation)

					rep = nirep
					code = nirep.Sequence().Clone()

					continue
				}
			}

			ctx.Set(0, method.Call(mrb, recv))
			mrb.callinfoPop()
		case op.Enter:
			_ = code.ReadW()
			ctx.SetSequenceCursor(code.Cursor())
		case op.EQ:
			a := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			val1 := ctx.Get(int(a))
			val2 := ctx.Get(int(a) + 1)

			ctx.Set(int(a), val1 == val2)
		case op.LT, op.LE, op.GT, op.GE:
			a := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			res, err := opCompare(ctx, int(a), int(a)+1, opCode)
			if err != nil {
				return nil, err
			}

			ctx.Set(int(a), res)
		case op.String:
			a := code.ReadB()
			b := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), rep.PoolValue(b))
		case op.Return:
			a := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ret := ctx.Get(int(a))
			if ctx.IsTop() {
				ctx.Set(int(rep.Locals()), ret)
				goto Stop
			}

			ctx.Set(0, ret)
			ci := mrb.callinfoPop()
			proc := ci.Proc()
			nirep, ok := proc.Body().(*insn.Representation)
			if !ok {
				return nil, ErrIRepNotFound
			}

			rep = nirep
			code = rep.Sequence().Clone()
			code.Seek(ci.GetSequnceCursor())
		case op.StrCat:
			a := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), fmt.Sprintf("%v%v", ctx.Get(int(a)), ctx.Get(int(a)+1)))
		case op.Class:
			a := code.ReadB()
			b := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

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
		case op.Method:
			a := code.ReadB()
			b := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			nirep := rep.Representation(b)
			nproc := mrb.procNew(nirep)

			ctx.Set(int(a), NewObjectValue(nproc))
		case op.Def:
			a := code.ReadB()
			b := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			targetClass := ctx.Get(int(a)).(RClass)
			methodProc := ctx.Get(int(a) + 1).(RProc)
			mid := rep.Symbol(b)

			method := newMethodFromProc(methodProc)
			mrb.defineMethodRaw(targetClass, mid, method)

			ctx.Set(int(a), mid)
		case op.TClass:
			a := code.ReadB()
			ctx.SetSequenceCursor(code.Cursor())

			ctx.Set(int(a), mrb.context.GetCallinfo().TargetClass())
		case op.Stop:
			ctx.SetSequenceCursor(code.Cursor())
			goto Stop
		default:
			return nil, fmt.Errorf("opcode %d not implemented", opCode)
		}
	}

Stop:
	return ctx.Get(int(rep.Locals())), nil
}

func (state *State) callinfoPush(pushStack int, cci uint8, targetClass RClass, proc RProc, block RProc, mid Symbol, argc uint16) *callinfo {
	ctx := state.context
	prevCi := ctx.GetCallinfo()

	if prevCi == nil {
		prevCi = &callinfo{}
	}

	callinfo := &callinfo{
		methodId:    mid,
		stackOffset: prevCi.stackOffset + pushStack,
		numArgs:     int(argc & 0xf),
		targetClass: targetClass,
		proc:        proc,
	}
	ctx.callinfo.Push(callinfo)

	return callinfo
}

func (state *State) callinfoPop() *callinfo {
	ctx := state.context
	ctx.callinfo.Pop()

	return ctx.GetCallinfo()
}

func opCompare(ctx *context, a int, b int, code op.Code) (bool, error) {
	val1, ok := toFloat64(ctx.Get(a))
	if !ok {
		return false, ErrNotPrimitiveTypeComparison
	}

	val2, ok := toFloat64(ctx.Get(b))
	if !ok {
		return false, ErrNotPrimitiveTypeComparison
	}

	switch code {
	case op.LT:
		return val1 < val2, nil
	case op.LE:
		return val1 <= val2, nil
	case op.GT:
		return val1 > val2, nil
	case op.GE:
		return val1 >= val2, nil
	}

	return false, ErrNotPrimitiveTypeComparison
}

func toFloat64(val Value) (float64, bool) {
	switch v := val.(type) {
	case int:
		return float64(v), true
	case float64:
		return v, true
	}

	return 0, false
}
