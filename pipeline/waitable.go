package pipeline

import "context"

type Waitable[T any] struct {
	data  T
	ready chan struct{}
}

func NewWaitable[T any]() Waitable[T] {
	return Waitable[T]{
		ready: make(chan struct{}),
	}
}

func (w *Waitable[T]) Set(data T) {
	w.data = data
	close(w.ready)
}

func (w *Waitable[T]) Get(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()

	case <-w.ready:
		return w.data, nil
	}
}

func (w *Waitable[T]) Ready() chan struct{} {
	return w.ready
}

func (w *Waitable[T]) Wait() {
	<-w.ready
}

func (w *Waitable[T]) Reset() {
	w.ready = make(chan struct{})

	var zero T
	w.data = zero
}
