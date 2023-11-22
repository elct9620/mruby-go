package stack

type Stack[T any] struct {
	container []T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{container: []T{}}
}

func (s *Stack[T]) Push(v T) {
	s.container = append(s.container, v)
}

func (s *Stack[T]) Pop() T {
	v := s.container[len(s.container)-1]
	s.container = s.container[:len(s.container)-1]
	return v
}

func (s *Stack[T]) Peek() T {
	return s.container[len(s.container)-1]
}

func (s *Stack[T]) Len() int {
	return len(s.container)
}

func (s *Stack[T]) Get(idx int) T {
	return s.container[idx]
}

func (s *Stack[T]) Set(idx int, v T) {
	s.container[idx] = v
}
