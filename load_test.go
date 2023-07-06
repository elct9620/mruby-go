package mruby_test

import (
	"bytes"
	"testing"

	"github.com/elct9620/mruby-go"
	"github.com/google/go-cmp/cmp"
)

func Test_ReadRiteHeader(t *testing.T) {
	expected := `#<RITE id="RITE" version="03.00" size="65" compiler="MATZ#0000">`

	buffer := bytes.NewBuffer(mustCompile(""))
	rite, err := mruby.ReadRiteHeader(buffer)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(expected, rite.String()) {
		t.Fatal("RITE header mismatched", cmp.Diff(expected, rite.String()))
	}
}
