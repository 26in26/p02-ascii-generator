package pipeline

import (
	"github.com/26in26/p02-ascii-generator/image"
)

type Stage interface {
	Process(ctx *FrameContext) error
}
type Pipeline struct {
	Stages []Stage
}

func New(process ...Stage) *Pipeline {
	return &Pipeline{Stages: process}
}

func (p *Pipeline) Run(img *image.Buffer) error {
	ctx := NewFrameContext(img)
	for _, stage := range p.Stages {
		err := stage.Process(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
