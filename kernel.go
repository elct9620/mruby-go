package mruby

func initKernel(mrb *State) {
	mrb.KernelModule = mrb.DefineModule("Kernel")
}
