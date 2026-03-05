package grayscale

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/utils"

	"github.com/26in26/p02-ascii-generator/pipeline"
)

type GrayScaleConnector struct{}

func NewGrayScaleConnector() pipeline.StageConnector[*image.RGBBuffer, *image.GrayBuffer] {
	return &GrayScaleConnector{}
}

func (c *GrayScaleConnector) Set(ctx context.Context, f *pipeline.Frame, val *image.GrayBuffer) error {
	f.GrayImage.Set(val)
	return nil
}

func (c *GrayScaleConnector) Get(ctx context.Context, f *pipeline.Frame) (*image.RGBBuffer, error) {
	srcImg, err := f.ResizedImage.Get(ctx)

	if err != nil {
		return nil, err
	}

	if srcImg == nil {
		return nil, utils.ErrBufferNotInitialized
	}

	return srcImg, nil
}
