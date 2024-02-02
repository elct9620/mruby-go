package mruby

func initCore(mrb *State) (err error) {
	err = initClass(mrb)
	if err != nil {
		return err
	}

	err = initObject(mrb)
	if err != nil {
		return err
	}
	initKernel(mrb)

	return nil
}
