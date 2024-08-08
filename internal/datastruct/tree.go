package datastruct

type Tree[K comparable, V any] struct {
	nodes map[K]*Tree[K, V]
	value *V
}

type ITree[K comparable, V any] interface {
	Get(keys ...K) (V, bool)
	Set(keys []K, v V)
	GetOrSet(keys []K, v V) (V, bool)
	// Iterator
}

var _ ITree[int, int] = &Tree[int, int]{}

func NewTree[K comparable, V any]() *Tree[K, V] {
	return &Tree[K, V]{
		nodes: map[K]*Tree[K, V]{},
	}
}

func (t *Tree[K, V]) Get(keys ...K) (V, bool) {
	if len(keys) == 0 {
		if t.value == nil {
			var ret V
			return ret, false
		}
		return *t.value, true
	}

	key := keys[0]
	node, ok := t.nodes[key]
	if !ok {
		var ret V
		return ret, false
	}
	return node.Get(keys[1:]...)
}
func (t *Tree[K, V]) Set(keys []K, value V) {
	if len(keys) == 0 {
		t.value = &value
		return
	}
	key := keys[0]
	node, ok := t.nodes[key]
	if !ok {
		node = NewTree[K, V]()
		t.nodes[key] = node
	}
	node.Set(keys[1:], value)
}
func (t *Tree[K, V]) GetOrSet(keys []K, value V) (V, bool) {
	v, exist := t.Get(keys...)
	if exist {
		return v, true
	}
	t.Set(keys, value)
	return value, false
}
