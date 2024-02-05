package stack_test

import (
	"testing"

	"github.com/elct9620/mruby-go/stack"
)

func Test_StackPush(t *testing.T) {
	s := stack.New[int](10)
	s.Push(1)
	s.Push(2)
	s.Push(3)

	if s.Peek() != 3 {
		t.Errorf("Expected 3, got %d", s.Peek())
	}
}

func Test_StackPushExpand(t *testing.T) {
	s := stack.New[int](2)
	s.Push(1)
	s.Push(2)
	s.Push(3)

	if s.Peek() != 3 {
		t.Errorf("Expected 3, got %d", s.Peek())
	}
}

func Test_StackPop(t *testing.T) {
	s := stack.New[int](10)
	s.Push(1)
	s.Push(2)

	if actual := s.Pop(); actual != 2 {
		t.Errorf("Expected 2, got %d", actual)
	}
}

func Test_StackPopOutOfBounds(t *testing.T) {
	s := stack.New[int](10)
	s.Push(1)
	s.Push(2)

	s.Pop()
	s.Pop()

	if actual := s.Pop(); actual != 0 {
		t.Errorf("Expected 0, got %d", actual)
	}
}

func Test_StackGet(t *testing.T) {
	s := stack.New[int](10)
	s.Push(1)
	s.Push(2)

	if actual := s.Get(0); actual != 1 {
		t.Errorf("Expected 1, got %d", actual)
	}
}

func Test_StackGetOutOfBounds(t *testing.T) {
	s := stack.New[int](2)
	s.Push(1)
	s.Push(2)

	if actual := s.Get(3); actual != 0 {
		t.Errorf("Expected 0, got %d", actual)
	}

	if actual := s.Get(-1); actual != 0 {
		t.Errorf("Expected 0, got %d", actual)
	}
}

func Test_StackSet(t *testing.T) {
	s := stack.New[int](10)
	s.Push(1)
	s.Push(2)

	s.Set(0, 3)

	if actual := s.Get(0); actual != 3 {
		t.Errorf("Expected 3, got %d", actual)
	}
}

func Test_StackSetOutOfBounds(t *testing.T) {
	s := stack.New[int](2)
	s.Push(1)
	s.Push(2)

	s.Set(3, 3)

	if actual := s.Get(3); actual != 0 {
		t.Errorf("Expected 0, got %d", actual)
	}

	s.Set(-1, 3)

	if actual := s.Get(-1); actual != 0 {
		t.Errorf("Expected 0, got %d", actual)
	}
}

func Test_StackSlice(t *testing.T) {
	s := stack.New[int](10)
	s.Push(1)
	s.Push(2)
	s.Push(3)

	if actual := s.Slice(0, 2); len(actual) != 2 {
		t.Errorf("Expected 2, got %d", len(actual))
	}
}

func Test_StackSliceOutOfBounds(t *testing.T) {
	s := stack.New[int](2)
	s.Push(1)
	s.Push(2)

	if actual := s.Slice(-1, 3); len(actual) != 2 {
		t.Errorf("Expected 2, got %d", len(actual))
	}
}
