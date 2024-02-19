package mruby

var _ RException = &Exception{}

type RException interface {
	RBasic
	error
	SetMessage(string)
}

type Exception struct {
	object
	message string
}

func (e *Exception) Class() RClass {
	return nil
}

func (e *Exception) Error() string {
	return e.message
}

func (e *Exception) SetMessage(message string) {
	e.message = message
}

func (mrb *State) ExceptionNewString(class RClass, message string) RException {
	exc := mrb.AllocException(class)
	exc.SetMessage(message)

	return exc
}

func (mrb *State) ExceptionMessageSet(exc RException, message string) {
	exc.SetMessage(message)
}
