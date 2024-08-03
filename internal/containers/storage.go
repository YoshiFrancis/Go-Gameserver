package containers

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

func (s *Storage[K, V]) Delete(key K) {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.data, key)
}

func (s *Storage[K, V]) Keys() []K {
	keys := make([]K, 0)
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}

func (s *Storage[K, V]) Values() []V {
	values := make([]V, 0)
	for _, value := range s.data {
		values = append(values, value)
	}

	return values
}

func (s *Storage[K, V]) Items() map[K]V {
	return s.data
}
