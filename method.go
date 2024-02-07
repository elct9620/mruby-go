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
	return m.Function(mrb, self)
}
