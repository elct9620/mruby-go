package insn_test

import (
	"bytes"
	"testing"

	"github.com/elct9620/mruby-go/insn"
)

type mockState struct {
	symCount uint32
	symTable map[string]uint32
}

func newMockState() *mockState {
	return &mockState{
		symTable: make(map[string]uint32),
	}
}

func (m *mockState) Intern(s string) uint32 {
	if _, ok := m.symTable[s]; !ok {
		m.symTable[s] = m.symCount + 1
	}

	return m.symTable[s]
}

func (m *mockState) Symbol(i uint32) string {
	for k, v := range m.symTable {
		if v == i {
			return k
		}
	}

	return ""
}

func Test_Representation_Sequence(t *testing.T) {
	buffer := bytes.NewReader([]byte{
		0x00, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
		0x08, 0x01, 0x38, 0x01, 0x69, 0x00, 0x00, 0x00, 0x00,
	})

	mrb := newMockState()
	reader := insn.NewBinaryReader(buffer)
	rep, err := insn.NewRepresentation(mrb, reader)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if rep.Sequence() == nil {
		t.Errorf("expected sequence to be not nil")
	}
}

func Test_Representation_Symbol(t *testing.T) {
	buffer := bytes.NewReader([]byte{
		0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x10, 0x01, 0x00, 0x38,
		0x01, 0x69, 0x00, 0x00, 0x00, 0x01, 0x00, 0x03, 0x6e, 0x65, 0x77, 0x00, 0x45, 0x4e, 0x44, 0x00,
	})

	mrb := newMockState()
	reader := insn.NewBinaryReader(buffer)
	rep, err := insn.NewRepresentation(mrb, reader)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if rep.Symbol(0) != 1 {
		t.Errorf("expected symbol to be 1 but got %v", rep.Symbol(0))
	}

	if mrb.Symbol(rep.Symbol(0)) != "new" {
		t.Errorf("expected symbol to be 'new' but got %v", mrb.Symbol(rep.Symbol(0)))
	}
}

func Test_Representation_PoolValue(t *testing.T) {
	buffer := bytes.NewReader([]byte{
		0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x51, 0x01, 0x00, 0x38,
		0x01, 0x69, 0x00, 0x01, 0x00, 0x00, 0x0c, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72,
		0x6c, 0x64, 0x21, 0x00, 0x00, 0x00, 0x45, 0x4e, 0x44, 0x00, 0x00, 0x00, 0x00, 0x08,
	})

	mrb := newMockState()
	reader := insn.NewBinaryReader(buffer)
	rep, err := insn.NewRepresentation(mrb, reader)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if rep.PoolValue(0) != "Hello World!" {
		t.Errorf("expected pool value to be 'Hello World!' but got %v", rep.PoolValue(0))
	}
}
