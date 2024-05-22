package mruby

func initArray(mrb *State) (err error) {
	mrb.ArrayClass, err = mrb.DefineClassId(_ArrayClass(mrb), mrb.ObjectClass)
	if err != nil {
		return err
	}

	mrb.DefineMethodId(mrb.ArrayClass, _join(mrb), arrayJoin)

	return nil
}

func arrayJoin(mrb *State, self Value) Value {
	argv := mrb.GetArgv()

	elem, ok := self.([]Value)
	if !ok {
		return nil
	}

	ret := ""
	elemLen := len(elem)
	for i, v := range elem {
		v, ok := v.(string)
		if !ok {
			continue
		}

		ret += v
		hasNext := elemLen > i+1
		if hasNext {
			ret += argv[0].(string)
		}
	}

	return ret
}
