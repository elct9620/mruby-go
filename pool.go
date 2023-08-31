package mruby

const (
	poolTypeString uint8 = iota
	poolTypeInt32
	poolTypeStaticString
	poolTypeInt64
	poolTypeFloat  = 5
	poolTypeBigInt = 7
)
