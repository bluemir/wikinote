package datastruct

import "container/list"

type Queue[T any] struct {
	items *list.List
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: list.New(),
	}
}
func (q *Queue[T]) Add(v T) {
	q.items.PushBack(v)
}
func (q *Queue[T]) Front() T {
	return q.items.Front().Value.(T)
}
func (q *Queue[T]) Pop() T {
	v := q.items.Front().Value
	q.items.Remove(q.items.Front())
	return v.(T)
}
func (q *Queue[T]) Len() int {
	return q.items.Len()
}
