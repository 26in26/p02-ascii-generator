package pipeline

import (
	"context"
)

type DataType int

const (
	DataRaw DataType = iota
	DataResized
	DataGray
	DataGradient
	DataAscii
)

type Stage interface {
	Name() string
	Requires() []DataType
	Provides() DataType
	Process(ctx context.Context, frame *Frame) error
}

type StageConnector[I any, O any] interface {
	Get(ctx context.Context, f *Frame) (I, error)
	Set(ctx context.Context, f *Frame, val O) error
}

type Kernal[I any, O any] func(ctx context.Context, input I) (O, error)

type BaseStage[I any, O any] struct {
	name     string
	requires []DataType
	provides DataType

	Connector StageConnector[I, O]
	Logic     Kernal[I, O]
}

func (s *BaseStage[I, O]) Name() string {
	return s.name
}

func (s *BaseStage[I, O]) Requires() []DataType {
	return s.requires
}

func (s *BaseStage[I, O]) Provides() DataType {
	return s.provides
}

func (s *BaseStage[I, O]) Process(ctx context.Context, frame *Frame) error {
	input, err := s.Connector.Get(ctx, frame)
	if err != nil {
		return err
	}

	output, err := s.Logic(ctx, input)
	if err != nil {
		return err
	}

	return s.Connector.Set(ctx, frame, output)
}

func NewBaseStage[I any, O any](
	name string,
	requires []DataType,
	provides DataType,
	connector StageConnector[I, O],
	Logic Kernal[I, O],
) *BaseStage[I, O] {

	return &BaseStage[I, O]{
		name:      name,
		requires:  requires,
		provides:  provides,
		Connector: connector,
		Logic:     Logic,
	}
}
