package storage

import "sync"

type Storage[K comparable, V any] struct {
	data map[K]V
	mux  sync.Mutex
}

func NewStorage[K comparable, V any]() *Storage[K, V] {
	return &Storage[K, V]{
		data: make(map[K]V),
		mux:  sync.Mutex{},
	}
}

func (s *Storage[K, V]) Get(key K) (value V, ok bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	value, ok = s.data[key]
	return value, ok
}

func (s *Storage[K, V]) Set(key K, value V) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.data[key] = value
}
