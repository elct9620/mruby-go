package mruby

type instructionSequence struct {
	binary []byte
	cursor int
}

func newInstructionSequence(bytes []byte) *instructionSequence {
	return &instructionSequence{binary: bytes, cursor: 0}
}

func (iseq *instructionSequence) ReadCode() code {
	return iseq.ReadB()
}

func (iseq *instructionSequence) ReadB() byte {
	b := iseq.binary[iseq.cursor]
	iseq.cursor++
	return b
}

func (iseq *instructionSequence) ReadS() []byte {
	s := iseq.binary[iseq.cursor : iseq.cursor+2]
	iseq.cursor += 2
	return s
}

func (iseq *instructionSequence) ReadW() []byte {
	w := iseq.binary[iseq.cursor : iseq.cursor+4]
	iseq.cursor += 4
	return w
}
