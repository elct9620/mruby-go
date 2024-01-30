package mruby

const nullSymbolLength = 0xffff

type iRepReaderFn func(*State, *iRep, *Reader) error

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
	iSeq      *Code
	poolValue []Value
	syms      []Symbol
}

func newIRep(mrb *State, r *Reader) (*iRep, error) {
	iRep := &iRep{}

	for _, loader := range iRepReaders {
		err := loader(mrb, iRep, r)
		if err != nil {
			return nil, err
		}
	}

	return iRep, nil
}

func readIRepHeader(mrb *State, ir *iRep, r *Reader) (err error) {
	fields := []any{
		&ir.nLocals,
		&ir.nRegs,
		&ir.rLen,
		&ir.cLen,
		&ir.iLen,
	}

	for _, field := range fields {
		err = r.ReadAs(field)
		if err != nil {
			return err
		}
	}

	return nil
}

func readISeq(mrb *State, ir *iRep, r *Reader) error {
	binary := make([]byte, ir.iLen)
	err := r.ReadAs(binary)
	if err != nil {
		return err
	}

	ir.iSeq = NewCode(binary)
	return nil
}

func readSyms(mrb *State, ir *iRep, r *Reader) error {
	err := r.ReadAs(&ir.sLen)
	if err != nil {
		return err
	}

	ir.syms = make([]Symbol, ir.sLen)

	for i := uint16(0); i < ir.sLen; i++ {
		strLen, err := r.ReadUint16()
		if err != nil {
			return err
		}

		if strLen == nullSymbolLength {
			ir.syms[i] = 0
			continue
		}

		str, err := r.ReadString(int(strLen) + 1)
		if err != nil {
			return err
		}

		ir.syms[i] = mrb.Intern(string(str))
	}

	return nil
}
