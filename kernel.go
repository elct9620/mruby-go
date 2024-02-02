package mruby

func initKernel(mrb *State) (err error) {
	mrb.KernelModule = mrb.DefineModuleId(_Kernel(mrb))

	return
}
