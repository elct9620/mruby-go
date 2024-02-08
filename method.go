package mruby

type Function func(*State, Value) Value

type Method interface {
	Call(*State, Value) Value
}

var _ Method = &method{}

type method struct {
	Function
	RProc
}

func (m *method) Call(mrb *State, self Value) Value {
	if m.RProc != nil {
		ret, err := mrb.VmRun(m.RProc, self)
		if err != nil {
			panic(err)
		}

		return ret
	}

	return m.Function(mrb, self)
}
