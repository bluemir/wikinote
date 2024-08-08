package datastruct

import (
	"sync"
)

type null struct{}

type Set[T comparable] struct {
	m map[T]null
	sync.RWMutex
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		m: map[T]null{},
	}
}

// Add add
func (s *Set[T]) Add(item T) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = null{}
}

// Remove deletes the specified item from the map
func (s *Set[T]) Remove(item T) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

// Has looks for the existence of an item
func (s *Set[T]) Has(item T) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

// Len returns the number of items in a set.
func (s *Set[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.List())
}

// Clear removes all items from the set
func (s *Set[T]) Clear() {
	s.Lock()
	defer s.Unlock()

	s.m = map[T]null{}
}

// IsEmpty checks for emptiness
func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Set returns a slice of all items
func (s *Set[T]) List() []T {
	s.RLock()
	defer s.RUnlock()

	list := make([]T, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
