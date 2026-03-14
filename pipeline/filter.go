package pipeline

import "context"

// Filter is a processing unit which modify unowned data
type Filter[I, O any] interface {
	Name() string
	Apply(ctx context.Context, input I) (O, error)
}

type BaseFilter[I, O any] struct {
	name string
}

func NewBaseFilter[I, O any](name string) *BaseFilter[I, O] {
	return &BaseFilter[I, O]{
		name: name,
	}
}

func (s *BaseFilter[I, O]) Name() string {
	return s.name
}

func (s *BaseFilter[I, O]) Apply(ctx context.Context, input I) (O, error) {
	panic("not implemented")
}
