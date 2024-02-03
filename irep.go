package mruby

import "github.com/elct9620/mruby-go/insn"

type iRepReaderFn func(*State, *iRep, insn.Reader) error

var iRepReaders = []iRepReaderFn{
	readIRepHeader,
	readISeq,
	readPoolValues,
	readSyms,
}

type iRep struct {
	nLocals   uint16
	nRegs     uint16
	rLen      uint16
	cLen      uint16
	iLen      uint32
	pLen      uint16
	sLen      uint16
	iSeq      *insn.Sequence
	poolValue []Value
	syms      []Symbol
}

func newIRep(mrb *State, r insn.Reader) (*iRep, error) {
	iRep := &iRep{}

	for _, loader := range iRepReaders {
		err := loader(mrb, iRep, r)
		if err != nil {
			return nil, err
		}
	}

	return iRep, nil
}

func readIRepHeader(mrb *State, ir *iRep, r insn.Reader) (err error) {
	fields := []any{
		&ir.nLocals,
		&ir.nRegs,
		&ir.rLen,
		&ir.cLen,
		&ir.iLen,
	}

	for _, field := range fields {
		err = r.As(field)
		if err != nil {
			return err
		}
	}

	return nil
}

func readISeq(mrb *State, ir *iRep, r insn.Reader) error {
	binary := make([]byte, ir.iLen)
	err := r.As(binary)
	if err != nil {
		return err
	}

	ir.iSeq = insn.NewSequence(binary)
	return nil
}

func readSyms(mrb *State, ir *iRep, r insn.Reader) error {
	err := r.As(&ir.sLen)
	if err != nil {
		return err
	}

	ir.syms = make([]Symbol, ir.sLen)

	for i := uint16(0); i < ir.sLen; i++ {
		symbol, err := r.String()
		if err != nil {
			return err
		}

		if len(symbol) == 0 {
			ir.syms[i] = 0
			continue
		}

		ir.syms[i] = mrb.Intern(symbol)
	}

	return nil
}
