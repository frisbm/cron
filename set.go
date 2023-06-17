package cron

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"sync"
)

type Set[T constraints.Ordered] struct {
	mu    sync.RWMutex
	items map[T]struct{}
}

func NewSet[T constraints.Ordered](items ...T) *Set[T] {
	set := &Set[T]{
		items: make(map[T]struct{}, len(items)),
	}
	set.Add(items...)
	return set
}

func (s *Set[T]) Add(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, item := range items {
		s.items[item] = struct{}{}
	}
}

func (s *Set[T]) Values() []T {
	values := maps.Keys(s.items)
	slices.Sort(values)
	return values
}
