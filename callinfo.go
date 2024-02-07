package mruby

type callinfo struct {
	numArgs     int
	methodId    Symbol
	stackOffset int
	targetClass RClass
}
