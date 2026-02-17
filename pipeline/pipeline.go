package pipeline

import (
	"github.com/26in26/p02-ascii-generator/image"
)

type Stage interface {
	Process(ctx *FrameContext)
}
type Pipeline struct {
	Stages []Stage
}

func New(process ...Stage) *Pipeline {
	return &Pipeline{Stages: process}
}

func (p *Pipeline) Run(img *image.Buffer) {
	ctx := NewFrameContext(img)
	for _, stage := range p.Stages {
		stage.Process(ctx)
	}
}
