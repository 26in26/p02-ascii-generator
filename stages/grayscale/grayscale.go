package grayscale

import (
	"fmt"

	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/utils"
)

type grayscaleStage struct{}

func (s *grayscaleStage) Process(ctx *pipeline.FrameContext) error {
	if ctx.WorkingImage == nil {
		return fmt.Errorf("grayscale stage: %w", utils.ErrBufferNotInitialized)
	}

	var err error
	ctx.GrayImage, err = ctx.WorkingImage.ToGray()

	return err
}

func NewGrayscaleStage() *grayscaleStage {
	return &grayscaleStage{}
}
