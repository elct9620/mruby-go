package mruby

type Function func(*State, Value) Value

type Method interface {
	Call(*State, Value) Value
	Proc() RProc
	IsProc() bool
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

func (m *method) Proc() RProc {
	return m.RProc
}

func (m *method) IsProc() bool {
	return m.RProc != nil
}
