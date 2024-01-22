package mruby

func initCore(mrb *State) (err error) {
	initClass(mrb)
	err = initObject(mrb)
	if err != nil {
		return err
	}
	initKernel(mrb)

	return nil
}
