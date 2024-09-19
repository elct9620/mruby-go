package mruby

import "fmt"

func objectInspect(mrb *State, self Value) Value {
	switch v := self.(type) {
	case *Object:
		return "Object"
	case *Class:
		name := mrb.ObjectInstanceVariableGet(v, _classname(mrb))
		return name
	case *Module:
		return "Module"
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
