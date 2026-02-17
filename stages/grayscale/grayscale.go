package grayscale

import (
	"github.com/26in26/p02-ascii-generator/pipeline"
)

type grayscaleStage struct{}

func (s *grayscaleStage) Process(ctx *pipeline.FrameContext) {
	if ctx.WorkingImage == nil {
		panic("grayscale stage must receive non-nil image buffer")
	}

	ctx.GrayImage = ctx.WorkingImage.ToGray()
}

func NewGrayscaleStage() *grayscaleStage {
	return &grayscaleStage{}
}
