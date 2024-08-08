package functional

func Map[In any, Out any](in []In, fn func(In) Out) []Out {
	out := []Out{}

	for _, v := range in {
		out = append(out, fn(v))
	}

	return out
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
