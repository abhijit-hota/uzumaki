package main

import "sync"

type SafeSet[T comparable] struct {
	values map[T]bool
	mu     sync.Mutex
}

func NewSafeSet[T comparable]() *SafeSet[T] {
	return &SafeSet[T]{
		values: make(map[T]bool),
	}
}

func (s *SafeSet[T]) Add(v T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, has := s.values[v]; has {
		return false
	}

	s.values[v] = true
	return true
}

func (s *SafeSet[T]) Values() []T {
	s.mu.Lock()
	defer s.mu.Unlock()

	values := make([]T, 0, len(s.values))
	for k := range s.values {
		values = append(values, k)
	}
	return values
}
