package builder

type Option[T any] func(*T)

func Build[T any](newT func() *T, opts ...Option[T]) *T {
	obj := newT() // create default
	for _, opt := range opts {
		opt(obj) // apply options
	}
	return obj
}
