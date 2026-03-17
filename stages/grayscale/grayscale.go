package grayscale

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
)

type GrayscaleStage struct {
	pipeline.BaseStage[*image.RGBBuffer, *image.GrayBuffer]
	pool *grayImagePool
}

func NewGrayscaleStage() pipeline.Stage[*image.RGBBuffer, *image.GrayBuffer] {
	return &GrayscaleStage{
		BaseStage: *pipeline.NewBaseStage[*image.RGBBuffer, *image.GrayBuffer]("Grayscale"),
		pool:      NewGrayImagePool(),
	}
}

func (s *GrayscaleStage) Kernal(ctx context.Context, input *image.RGBBuffer) (*image.GrayBuffer, error) {
	grayBuffer, err := s.pool.Get(size{w: input.Width, h: input.Height})
	if err != nil {
		return nil, err
	}

	return input.ToGray(grayBuffer)
}

func (s *GrayscaleStage) Release(input *image.GrayBuffer) {
	if input == nil {
		return
	}

	s.pool.Put(size{w: input.Width, h: input.Height}, input)
}
