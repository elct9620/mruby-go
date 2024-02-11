package mruby

func NewObjectValue(v any) Value {
	return v.(RBasic)
}

func Bool(v Value) bool {
	ret, ok := v.(bool)

	if !ok {
		return false
	}

	return ret
}

func Test(v Value) bool {
	return Bool(v)
}
