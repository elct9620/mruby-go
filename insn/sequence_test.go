package insn_test

import (
	"testing"

	"github.com/elct9620/mruby-go/insn"
	"github.com/elct9620/mruby-go/op"
)

func Test_Operation(t *testing.T) {
	iseq := insn.NewSequence([]byte{
		op.Return,
	})

	code := iseq.Operation()

	if code != op.Return {
		t.Errorf("Expected op.Return, got %v", code)
	}
}

func Test_ReadB(t *testing.T) {
	iseq := insn.NewSequence([]byte{
		0x01,
	})

	b := iseq.ReadB()

	if b != 0x01 {
		t.Errorf("Expected 0x01, got %v", b)
	}
}

func Test_ReadS(t *testing.T) {
	iseq := insn.NewSequence([]byte{
		0x01, 0x02,
	})

	s := iseq.ReadS()

	if len(s) != 2 {
		t.Errorf("Expected reads 2 bytes, got %v", len(s))
	}

	if s[0] != 0x01 || s[1] != 0x02 {
		t.Errorf("Expected 0x01 0x02, got %v %v", s[0], s[1])
	}
}

func Test_ReadW(t *testing.T) {
	iseq := insn.NewSequence([]byte{
		0x01, 0x02, 0x03,
	})

	w := iseq.ReadW()

	if len(w) != 3 {
		t.Errorf("Expected reads 3 bytes, got %v", len(w))
	}

	if w[0] != 0x01 || w[1] != 0x02 || w[2] != 0x03 {
		t.Errorf("Expected 0x01 0x02 0x03, got %v %v %v", w[0], w[1], w[2])
	}
}

func Test_Cursor(t *testing.T) {
	iseq := insn.NewSequence([]byte{
		0x01, 0x02,
	})

	iseq.ReadB()

	cursor := iseq.Cursor()

	if cursor != 1 {
		t.Errorf("Expected cursor 1, got %v", cursor)
	}
}

func Test_Clone(t *testing.T) {
	iseq := insn.NewSequence([]byte{
		0x01, 0x02,
	})

	iseq.ReadB()

	clone := iseq.Clone()
	b := clone.ReadB()

	if b != 0x02 {
		t.Errorf("Expected 0x02, got %v", b)
	}
}
