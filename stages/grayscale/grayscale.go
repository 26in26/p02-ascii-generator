package grayscale

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
)

func Grayscale(ctx context.Context, input *image.RGBBuffer) (*image.GrayBuffer, error) {
	return input.ToGray()
}

func NewGrayscaleStage() pipeline.Stage {
	return pipeline.NewBaseStage("grayscale", []pipeline.DataType{pipeline.DataResized}, pipeline.DataGray,
		NewGrayScaleConnector(),
		Grayscale,
	)
}
