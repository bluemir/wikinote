package datastruct

import "sync"

type Map[K comparable, V any] struct {
	internal sync.Map
}

func (m *Map[K, V]) Set(k K, v V) {
	m.internal.Store(k, v)
}
func (m *Map[K, V]) Get(k K) (V, bool) {
	v, ok := m.internal.Load(k)
	if !ok {
		var ret V
		return ret, false
	}

	val, ok := v.(V)
	return val, ok
}
func (m *Map[K, V]) GetOrSet(k K, v V) (V, bool) {
	val, ok := m.internal.LoadOrStore(k, v)
	if !ok {
		return val.(V), false
	}

	ret, ok := val.(V)
	return ret, ok
}
func (m *Map[K, V]) Range(fn func(k K, v V) bool) {
	m.internal.Range(func(key, value any) bool {
		return fn(key.(K), value.(V))
	})
}
