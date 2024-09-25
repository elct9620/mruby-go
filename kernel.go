package mruby

import "fmt"

func (mrb *State) Inspect(obj Value) string {
	ret, ok := objectInspect(mrb, obj).(string)
	if !ok {
		return ""
	}

	return ret
}

func objectInspect(mrb *State, self Value) Value {
	switch v := self.(type) {
	case *Object:
		name := mrb.ObjectInstanceVariableGet(v.Class(), _classname(mrb))
		return name
	case RClass:
		name := mrb.ObjectInstanceVariableGet(v, _classname(mrb))
		return name
	case string:
		return fmt.Sprintf("%q", v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	default:
		return nil
	}
}

func objectPuts(mrb *State, recv Value) Value {
	args := mrb.GetArgv()
	fmt.Println(args...)
	return args[0]
}

func funcRaise(mrb *State, recv Value) Value {
	args := mrb.GetArgv()
	mrb.Raise(nil, fmt.Sprint(args...))

	return nil
}

func initKernel(mrb *State) (err error) {
	mrb.KernelModule = mrb.DefineModuleId(_Kernel(mrb))

	mrb.DefineMethodId(mrb.KernelModule, _inspect(mrb), objectInspect)
	mrb.DefineMethodId(mrb.ObjectClass, _raise(mrb), funcRaise)

	// defined in mruby-io
	putsSym := mrb.Intern("puts")
	mrb.DefineMethodId(mrb.ObjectClass, putsSym, objectPuts)

	return mrb.IncludeModule(mrb.ObjectClass, mrb.KernelModule)
}
