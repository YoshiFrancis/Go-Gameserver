package containers

import (
	"fmt"
	"sync"
)

type Queue[T any] struct {
	maxSize int
	data    []T
	mux     sync.Mutex
}

func NewQueue[T any](maxSize int) *Queue[T] {
	return &Queue[T]{
		maxSize: maxSize,
		data:    make([]T, 0),
		mux:     sync.Mutex{},
	}
}

func (q *Queue[T]) Items() []T {
	q.mux.Lock()
	defer q.mux.Unlock()
	return q.data
}

func (q *Queue[T]) Dequeue() T {
	q.mux.Lock()
	defer q.mux.Unlock()
	item := q.data[0]
	q.data = q.data[1:]
	return item
}

func (q *Queue[T]) Enqueue(item T) {
	q.mux.Lock()
	defer q.mux.Unlock()
	fmt.Println(len(q.data), " == ", q.maxSize)
	if len(q.data) == q.maxSize {
		q.data = q.data[1:] // dequeuing. I did not call Dequeue due to me having to unlocking mux then relock afterward
	}
	q.data = append(q.data, item)
}

func (q *Queue[T]) Front() T {
	q.mux.Lock()
	defer q.mux.Unlock()
	return q.data[0]
}

func (q *Queue[T]) Back() T {
	q.mux.Lock()
	defer q.mux.Unlock()
	return q.data[len(q.data)-1]
}
