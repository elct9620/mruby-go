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
