package mruby

type (
	value = any
	code  = uint8
)

type state struct {
}

func New() *state {
	return &state{}
}
