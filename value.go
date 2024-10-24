package mruby

func NewObjectValue(v any) Value {
	return v.(RBasic)
}

func Bool(v Value) bool {
	switch v := v.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}

func Test(v Value) bool {
	return Bool(v)
}

func ClassP(v Value) bool {
	switch v.(type) {
	case *Class:
		return true
	default:
		return false
	}
}
