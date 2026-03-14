package pipeline

import "context"

// Stage is a proccessing unit which allocates it's own data.
// A stage has an explicit release function which frees the output memory (Reducing GC load)
type Stage[I, O any] interface {
	Name() string
	Kernal(ctx context.Context, input I) (O, error)
	Release(O)
}

// A default stage with all the required functions
type BaseStage[I, O any] struct {
	name string
}

func NewBaseStage[I, O any](name string) *BaseStage[I, O] {
	return &BaseStage[I, O]{
		name: name,
	}
}

func (s *BaseStage[I, O]) Name() string {
	return s.name
}

func (s *BaseStage[I, O]) Kernal(ctx context.Context, input I) (O, error) {
	panic("not implemented")
}

func (s *BaseStage[I, O]) Release(O) {}
