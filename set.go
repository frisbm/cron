package cron

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type set[T constraints.Ordered] struct {
	items map[T]struct{}
}

func newSet[T constraints.Ordered](items ...T) *set[T] {
	s := &set[T]{
		items: make(map[T]struct{}, len(items)),
	}
	s.Add(items...)
	return s
}

func (s *set[T]) Add(items ...T) {
	for _, item := range items {
		s.items[item] = struct{}{}
	}
}

func (s *set[T]) Values() []T {
	values := maps.Keys(s.items)
	slices.Sort(values)
	return values
}
