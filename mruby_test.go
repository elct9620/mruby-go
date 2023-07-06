package mruby_test

import (
	"testing"

	"github.com/elct9620/mruby-go"
	"github.com/google/go-cmp/cmp"
)

func Test_Mrb_Header(t *testing.T) {
	expected := `#<RITE id="RITE" version="03.00" size="65" compiler="MATZ#0000">`

	mrb, err := mruby.NewFromString("")
	if err != nil {
		t.Fatal(err)
	}

	riteHeader := mrb.Header().String()
	if !cmp.Equal(expected, riteHeader) {
		t.Fatal("RITE header mismatched", cmp.Diff(expected, riteHeader))
	}
}
