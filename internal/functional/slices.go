package functional

func Map[In any, Out any](in []In, fn func(In) Out) []Out {
	out := make([]Out, 0, len(in))

	for _, v := range in {
		out = append(out, fn(v))
	}

	return out
}
func MapWithError[In any, Out any](in []In, fn func(In) (Out, error)) ([]Out, error) {
	out := make([]Out, 0, len(in))

	for _, v := range in {
		o, err := fn(v)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}

	return out, nil
}
func Contain[In comparable](in []In, v In) bool {
	return ContainWithFn(in, func(i In) bool { return i == v })
}
func ContainWithFn[In any](in []In, fn func(In) bool) bool {
	for _, v := range in {
		if fn(v) {
			return true
		}
	}

	return false
}
func Filter[T any](in []T, fn func(T) bool) []T {
	out := []T{}

	for _, v := range in {
		if fn(v) {
			out = append(out, v)
		}
	}

	return out
}
func ToLookupTable[Key comparable, Element any](in []Element, keyFn func(Element) Key) map[Key]Element {
	m := map[Key]Element{}

	for _, elem := range in {
		key := keyFn(elem)
		m[key] = elem
	}

	return m
}
func Some[T any](in []T, fn func(T) bool) bool {
	for _, v := range in {
		if fn(v) {
			return true
		}
	}
	return false
}
func All[T any](in []T, fn func(T) bool) bool {
	for _, v := range in {
		if !fn(v) {
			return false
		}
	}
	return true
}
func Reduce[In any, Out any](in []In, fn func(accumulator *Out, v In) Out, init Out) Out {
	out := init

	for _, v := range in {
		fn(&out, v)
	}

	return out
}
func Flat[T any](in [][]T) []T {
	out := []T{}

	for _, arr := range in {
		out = append(out, arr...)
	}
	return out
}

func First[In any](in []In, fn func(In) bool) *In {
	for _, v := range in {
		if fn(v) {
			return &v
		}
	}
	return nil
}

// arr = Filter(arr, func(v int) bool)
// arr = Map(arr, func(v int) float)
// ---- vs ----
// functional.From(arr).Filter(func(v int) bool).Map(func(v int) float).ToArray()
