package insn_test

import (
	"bytes"
	"testing"

	"github.com/elct9620/mruby-go/insn"
)

func Test_Reader_As(t *testing.T) {
	data := []byte{0x01}
	reader := insn.NewBinaryReader(bytes.NewReader(data))

	var i uint8
	if err := reader.As(&i); err != nil {
		t.Errorf("Read As failed: %v", err)
	}

	if i != 1 {
		t.Errorf("Read As failed: expected 1, got %v", i)
	}
}

func Test_Reader_Uint16(t *testing.T) {
	data := []byte{0x00, 0x02}
	reader := insn.NewBinaryReader(bytes.NewReader(data))

	i, err := reader.Uint16()
	if err != nil {
		t.Errorf("Read Uint16 failed: %v", err)
	}

	if i != 2 {
		t.Errorf("Uint16 failed: expected 2, got %v", i)
	}
}

func Test_Reader_Uint16_Invalid(t *testing.T) {
	data := []byte{0x00}
	reader := insn.NewBinaryReader(bytes.NewReader(data))

	_, err := reader.Uint16()
	if err == nil {
		t.Errorf("Read Uint16 failed: expected error, got nil")
	}
}

func Test_Reader_String(t *testing.T) {
	data := []byte{0x00, 0x03, 'a', 'b', 'c', 0x00}
	reader := insn.NewBinaryReader(bytes.NewReader(data))

	s, err := reader.String()
	if err != nil {
		t.Errorf("Read String failed: %v", err)
	}

	if s != "abc" {
		t.Errorf("Read String failed: expected abc, got %v", s)
	}
}

func Test_Reader_String_Length_Invalid(t *testing.T) {
	data := []byte{0x00}
	reader := insn.NewBinaryReader(bytes.NewReader(data))

	_, err := reader.String()
	if err == nil {
		t.Errorf("Read String failed: expected error, got nil")
	}
}

func Test_Reader_String_Invalid(t *testing.T) {
	data := []byte{0x00, 0x03, 'a', 'b'}
	reader := insn.NewBinaryReader(bytes.NewReader(data))

	_, err := reader.String()
	if err == nil {
		t.Errorf("Read String failed: expected error, got nil")
	}
}
