package pipeline

import "github.com/26in26/p02-ascii-generator/image"

type Stage interface {
	Process(img *image.Buffer) *image.Buffer
}

type Pipeline struct {
	stages []Stage
}

func New(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

func (p *Pipeline) Run(img *image.Buffer) *image.Buffer {
	out := img
	for _, stage := range p.stages {
		out = stage.Process(out)
	}
	return out
}
