package mruby

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

func initKernel(mrb *State) (err error) {
	mrb.KernelModule = mrb.DefineModuleId(_Kernel(mrb))

	mrb.DefineMethodId(mrb.KernelModule, _inspect(mrb), objectInspect)

	return mrb.IncludeModule(mrb.ObjectClass, mrb.KernelModule)
}
