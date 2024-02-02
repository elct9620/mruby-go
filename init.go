package mruby

func initCore(mrb *State) (err error) {
	if err = initClass(mrb); err != nil {
		return
	}

	if err = initObject(mrb); err != nil {
		return
	}

	if err = initKernel(mrb); err != nil {
		return
	}

	return
}
