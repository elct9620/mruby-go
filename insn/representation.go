package insn

type Symbol = uint32
type Value = any

type State interface {
	Intern(string) Symbol
}

type loadRepresentationFn func(State, *Representation, Reader) error

type Representation struct {
	nLocals   uint16
	nRegs     uint16
	rLen      uint16
	cLen      uint16
	iLen      uint32
	pLen      uint16
	sLen      uint16
	iSeq      *Sequence
	poolValue []Value
	syms      []Symbol
	reps      []*Representation
}

func NewRepresentation(mrb State, r Reader) (*Representation, error) {
	var rep = &Representation{}

	for _, fn := range []loadRepresentationFn{
		loadHeader,
		loadSequence,
		loadPool,
		loadSyms,
		loadReps,
	} {
		if err := fn(mrb, rep, r); err != nil {
			return nil, err
		}
	}

	return rep, nil
}

func (rep *Representation) Sequence() *Sequence {
	return rep.iSeq
}

func (rep *Representation) Symbol(i uint8) Symbol {
	return rep.syms[i]
}

func (rep *Representation) PoolValue(i uint8) Value {
	return rep.poolValue[i]
}

func (rep *Representation) Representation(i uint8) *Representation {
	return rep.reps[i]
}

func loadHeader(mrb State, rep *Representation, r Reader) error {
	fields := []any{
		&rep.nLocals,
		&rep.nRegs,
		&rep.rLen,
		&rep.cLen,
	}

	for _, field := range fields {
		if err := r.As(field); err != nil {
			return err
		}
	}

	return nil
}

func loadSequence(mrb State, rep *Representation, r Reader) error {
	if err := r.As(&rep.iLen); err != nil {
		return err
	}

	binary := make([]byte, rep.iLen)
	if err := r.As(binary); err != nil {
		return err
	}

	rep.iSeq = NewSequence(binary)
	return nil
}

func loadPool(mrb State, rep *Representation, r Reader) error {
	pLen, err := r.Uint16()
	if err != nil {
		return err
	}

	rep.poolValue = make([]any, pLen)

	for i := uint16(0); i < pLen; i++ {
		rep.poolValue[i], err = readPool(r)
		if err != nil {
			return err
		}

		rep.pLen = uint16(i + 1)
	}

	return nil
}

func loadSyms(mrb State, rep *Representation, r Reader) error {
	if err := r.As(&rep.sLen); err != nil {
		return err
	}

	rep.syms = make([]Symbol, rep.sLen)

	for i := uint16(0); i < rep.sLen; i++ {
		symbol, err := r.String()
		if err != nil {
			return err
		}

		if len(symbol) == 0 {
			rep.syms[i] = 0
			continue
		}

		rep.syms[i] = mrb.Intern(symbol)
	}

	return nil
}

func loadReps(mrb State, rep *Representation, r Reader) (err error) {
	rep.reps = make([]*Representation, rep.rLen)

	for i := uint16(0); i < rep.rLen; i++ {
		var size uint32
		if err = r.As(&size); err != nil {
			return
		}

		rep.reps[i], err = NewRepresentation(mrb, r)
		if err != nil {
			return
		}
	}

	return
}
