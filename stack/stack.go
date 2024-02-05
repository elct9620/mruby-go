package stack

type Stack[T any] struct {
	current  int
	end      int
	elements []T
}

func New[T any](size int) *Stack[T] {
	return &Stack[T]{
		current:  0,
		end:      size,
		elements: make([]T, size),
	}
}

func (s *Stack[T]) Push(v T) {
	s.current++

	if s.current == s.end {
		s.elements = append(s.elements, v)
		s.end++
		return
	}

	s.elements[s.current] = v
}

func (s *Stack[T]) Pop() T {
	var empty T
	if s.current == 0 {
		return empty
	}

	v := s.elements[s.current]
	s.elements[s.current] = empty
	s.current--

	return v
}

func (s *Stack[T]) Peek() T {
	return s.elements[s.current]
}
