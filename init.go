package mruby

func initCore(mrb *State) {
	initClass(mrb)
	initObject(mrb)
	initKernel(mrb)
}
