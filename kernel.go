package mruby

func objectInspect(mrb *State, self Value) Value {
	return nil
}

func initKernel(mrb *State) (err error) {
	mrb.KernelModule = mrb.DefineModuleId(_Kernel(mrb))

	mrb.DefineMethodId(mrb.KernelModule, _inspect(mrb), objectInspect)

	return mrb.IncludeModule(mrb.ObjectClass, mrb.KernelModule)
}
