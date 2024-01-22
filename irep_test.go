package mruby_test

import (
	"testing"

	"github.com/elct9620/mruby-go"
	"github.com/google/go-cmp/cmp"
)

func TestLoadInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1 - 2", -1},
		{"1 - 1", 0},
		{"2 - 1", 1},
		{"1 + 1", 2},
		{"1 + 2", 3},
		{"1 + 3", 4},
		{"2 + 3", 5},
		{"4 + 2", 6},
		{"4 + 3", 7},
	}

	mrb, err := mruby.New()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			res, err := mrb.LoadString(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !cmp.Equal(res, tc.expected) {
				t.Fatalf("expected %v, got %v", tc.expected, res)
			}
		})
	}
}
