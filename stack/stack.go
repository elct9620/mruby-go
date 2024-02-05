package stack

type Stack[T any] struct {
	current  int
	end      int
	elements []T
}

func New[T any](size int) *Stack[T] {
	return &Stack[T]{
		current:  -1,
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
	if s.current == -1 {
		var empty T
		return empty
	}

	return s.elements[s.current]
}

func (s *Stack[T]) Slice(start, end int) []T {
	if start < 0 {
		start = 0
	}

	if end > s.end {
		end = s.end
	}

	return s.elements[start:end]
}

func (s *Stack[T]) Get(idx int) T {
	if idx < 0 || idx > s.end {
		var empty T
		return empty
	}

	return s.elements[idx]
}

func (s *Stack[T]) Set(idx int, v T) {
	if idx < 0 || idx > s.end {
		return
	}

	s.elements[idx] = v
}
