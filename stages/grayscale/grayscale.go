package grayscale

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
)

type GrayscaleStage struct {
	pipeline.BaseStage[*image.RGBBuffer, *image.GrayBuffer]
}

func NewGrayscaleStage() pipeline.Stage[*image.RGBBuffer, *image.GrayBuffer] {
	return &GrayscaleStage{
		BaseStage: *pipeline.NewBaseStage[*image.RGBBuffer, *image.GrayBuffer]("Grayscale"),
	}
}

func (s *GrayscaleStage) Kernal(ctx context.Context, input *image.RGBBuffer) (*image.GrayBuffer, error) {
	return input.ToGray()
}
