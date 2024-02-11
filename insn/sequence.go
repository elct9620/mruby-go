package insn

import "github.com/elct9620/mruby-go/op"

type Sequence struct {
	code   []byte
	cursor int
}

func NewSequence(code []byte) *Sequence {
	return &Sequence{code: code, cursor: 0}
}

// Operation returns the next operation code in the sequence.
func (s *Sequence) Operation() op.Code {
	return s.ReadB()
}

// ReadB reads a byte from the sequence and advances the cursor.
func (s *Sequence) ReadB() byte {
	b := s.code[s.cursor]
	s.cursor++
	return b
}

// ReadS reads 2 bytes from the sequence and advances the cursor.
func (s *Sequence) ReadS() []byte {
	b := s.code[s.cursor : s.cursor+2]
	s.cursor += 2
	return b
}

// ReadW reads 3 bytes from the sequence and advances the cursor.
func (s *Sequence) ReadW() []byte {
	b := s.code[s.cursor : s.cursor+3]
	s.cursor += 3
	return b
}

// Cursor returns the current position of the cursor.
func (s *Sequence) Cursor() int {
	return s.cursor
}

// Seek sets the cursor to the given position.
func (s *Sequence) Seek(pos int) {
	s.cursor = pos
}

// Clone returns a new sequence with the same code and cursor.
func (s *Sequence) Clone() *Sequence {
	return &Sequence{code: s.code, cursor: s.cursor}
}
