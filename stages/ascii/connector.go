package ascii

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/utils"
)

type AsciiConnector struct{}

func NewAsciiConnector() pipeline.StageConnector[*AsciiInput, *image.AsciiBuffer] {
	return &AsciiConnector{}
}

func (c *AsciiConnector) Set(ctx context.Context, f *pipeline.Frame, val *image.AsciiBuffer) error {
	f.AsciiArt.Set(val)
	return nil
}

func (c *AsciiConnector) Get(ctx context.Context, f *pipeline.Frame) (*AsciiInput, error) {
	workingImg, err := f.ResizedImage.Get(ctx)

	if err != nil {
		return nil, err
	}

	if workingImg == nil {
		return nil, utils.ErrBufferNotInitialized
	}
	grayImg, err := f.GrayImage.Get(ctx)
	if err != nil {
		return nil, err
	}

	if grayImg == nil {
		return nil, utils.ErrBufferNotInitialized
	}

	gradient, err := f.GradientMap.Get(ctx)
	if err != nil {
		return nil, err
	}
	if gradient == nil {
		return nil, utils.ErrBufferNotInitialized
	}

	return &AsciiInput{
		workingImg: workingImg,
		grayImg:    grayImg,
		gradient:   gradient,
	}, nil
}
