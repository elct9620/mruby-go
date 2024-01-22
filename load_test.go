package mruby_test

import (
	"bytes"
	"testing"

	"github.com/elct9620/mruby-go"
	"github.com/google/go-cmp/cmp"
)

func Test_Mrb_Load(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{
		0x52, 0x49, 0x54, 0x45, 0x30, 0x33, 0x30, 0x30, 0x00, 0x00, 0x00, 0x66, 0x4d, 0x41, 0x54, 0x5a,
		0x30, 0x30, 0x30, 0x30, 0x49, 0x52, 0x45, 0x50, 0x00, 0x00, 0x00, 0x25, 0x30, 0x33, 0x30, 0x30,
		0x00, 0x00, 0x00, 0x19, 0x00, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
		0x08, 0x01, 0x38, 0x01, 0x69, 0x00, 0x00, 0x00, 0x00, 0x44, 0x42, 0x47, 0x00, 0x00, 0x00, 0x00,
		0x25, 0x00, 0x01, 0x00, 0x06, 0x61, 0x64, 0x64, 0x2e, 0x72, 0x62, 0x00, 0x00, 0x00, 0x13, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x02, 0x00, 0x03, 0x45, 0x4e,
		0x44, 0x00, 0x00, 0x00, 0x00, 0x08,
	})

	mrb, err := mruby.New()
	if err != nil {
		t.Fatal(err)
	}

	res, err := mrb.Load(buffer)
	if err != nil {
		t.Fatal(err)
	}

	expected := 2
	if !cmp.Equal(expected, res) {
		t.Fatal("return value mismatch", cmp.Diff(expected, res))
	}
}

func Test_Mrb_LoadString(t *testing.T) {
	mrb, err := mruby.New()
	if err != nil {
		t.Fatal(err)
	}

	res, err := mrb.LoadString("1 + 1")
	if err != nil {
		t.Fatal(err)
	}

	expected := 2
	if !cmp.Equal(expected, res) {
		t.Fatal("return value mismatch", cmp.Diff(expected, res))
	}
}
