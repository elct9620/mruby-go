package mruby

func NewObjectValue(v any) Value {
	return v.(RObject)
}
