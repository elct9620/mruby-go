package mruby_test

import (
	"bytes"
	"testing"

	"github.com/elct9620/mruby-go"
	"github.com/google/go-cmp/cmp"
)

func Test_Compile(t *testing.T) {
	expected := []byte{
		0x52, 0x49, 0x54, 0x45, 0x30, 0x33, 0x30, 0x30, 0x00, 0x00, 0x00, 0x41, 0x4d, 0x41, 0x54, 0x5a,
		0x30, 0x30, 0x30, 0x30, 0x49, 0x52, 0x45, 0x50, 0x00, 0x00, 0x00, 0x25, 0x30, 0x33, 0x30, 0x30,
		0x00, 0x00, 0x00, 0x19, 0x00, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
		0x08, 0x01, 0x38, 0x01, 0x69, 0x00, 0x00, 0x00, 0x00, 0x45, 0x4e, 0x44, 0x00, 0x00, 0x00, 0x00,
		0x08,
	}

	bin := mustCompile("1 + 1")
	if !cmp.Equal(bin, expected) {
		t.Fatal("compiled binary mismatched", cmp.Diff(expected, bin))
	}
}

func mustCompile(code string) []byte {
	bin, err := mruby.Compile(bytes.NewBufferString(code))
	if err != nil {
		panic(err)
	}

	return bin
}
