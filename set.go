package cron

import (
	"golang.org/x/exp/constraints"
)

type set[T constraints.Ordered] struct {
	items map[T]struct{}
}

func newSet[T constraints.Ordered](capacity int, items ...T) *set[T] {
	s := &set[T]{
		items: make(map[T]struct{}, capacity),
	}
	s.Add(items...)
	return s
}

func (s set[T]) Add(items ...T) {
	for _, item := range items {
		s.items[item] = struct{}{}
	}
}

func (s set[T]) Contains(key T) bool {
	_, ok := s.items[key]
	return ok
}
