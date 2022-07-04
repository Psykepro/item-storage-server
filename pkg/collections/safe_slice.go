package collections

import "sync"

type SafeSlice[T any] struct {
	sync.RWMutex
	slice []T
}

func NewSafeSlice[T any](size, capacity int) *SafeSlice[T] {
	return &SafeSlice[T]{
		slice: make([]T, size, capacity),
	}
}

func (s *SafeSlice[T]) Append(element T) {
	s.Lock()
	defer s.Unlock()
	
	s.slice = append(s.slice, element)
}

func (s *SafeSlice[T]) Set(index int, element T) {
	s.Lock()
	defer s.Unlock()

	s.slice[index] = element
}
func (s *SafeSlice[T]) Get(index int) T {
	s.RLock()
	defer s.RUnlock()

	return s.slice[index]
}

func (s *SafeSlice[T]) Remove(element T) {
	s.Lock()
	defer s.Unlock()

	s.slice = append(s.slice, element)
}

func (s *SafeSlice[T]) Slice() []T {
	s.RLock()
	defer s.RUnlock()

	return s.slice
}
