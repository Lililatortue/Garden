package types

import "errors"

func (s List[T]) IsEmpty() bool {
	return len(s) == 0
}

func (s List[T]) Push(node ...*T) {
	s = append(s, node...)
}

func (s List[T]) Pop() (*T, error) {
	if s.IsEmpty() {
		return nil, errors.New("no more nodes to pop")
	}
	last := s[len(s)-1]
	s = s[:len(s)-1]
	return last, nil
}

func (s List[T]) Shift() (*T, error) {
	if s.IsEmpty() {
		return nil, errors.New("no more nodes to shift")
	}
	first := s[0]
	s = s[1:]
	return first, nil
}
