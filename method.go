package mruby

type Function func(*State, Value) Value

type Method interface {
	Call(*State, Value) Value
}

var _ Method = &goMethod{}

type goMethod struct {
	Function
}

func (m *goMethod) Call(mrb *State, self Value) Value {
	return m.Function(mrb, self)
}
