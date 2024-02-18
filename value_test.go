package mruby_test

import (
	"testing"

	"github.com/elct9620/mruby-go"
)

func Test_Bool(t *testing.T) {
	tests := []struct {
		name     string
		actual   any
		expected bool
	}{
		{"true", true, true},
		{"false", false, false},
		{"nil", nil, false},
		{"int", 1, true},
		{"string", "string", true},
		{"object", &mruby.Object{}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test := test
			t.Parallel()

			if actual := mruby.Bool(test.actual); actual != test.expected {
				t.Errorf("expected %v, but got %v", test.expected, actual)
			}
		})
	}
}
