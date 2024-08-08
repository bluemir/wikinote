package datastruct

import "container/list"

type Queue[T any] struct {
	internal *list.List
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		internal: list.New(),
	}
}

func (q *Queue[T]) Add(v T) {
	q.internal.PushBack(v)
}
func (q *Queue[T]) Front() T {
	return q.internal.Front().Value.(T)
}
func (q *Queue[T]) Pop() T {
	v := q.internal.Front().Value
	q.internal.Remove(q.internal.Front())
	return v.(T)
}
func (q *Queue[T]) Len() int {
	return q.internal.Len()
}
