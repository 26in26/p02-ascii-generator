package resize

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/utils"
)

type ResizeConnector struct {
	targetWidth  int
	targetHeight int
}

func NewResizeConnector(w, h int) pipeline.StageConnector[*ResizeInput, *image.RGBBuffer] {
	return &ResizeConnector{
		targetWidth:  w,
		targetHeight: h,
	}
}

func (c *ResizeConnector) Set(ctx context.Context, f *pipeline.Frame, val *image.RGBBuffer) error {
	f.ResizedImage.Set(val)
	return nil
}

func (c *ResizeConnector) Get(ctx context.Context, f *pipeline.Frame) (*ResizeInput, error) {
	srcImg, err := f.Raw.Get(ctx)

	if err != nil {
		return nil, err
	}

	if srcImg == nil {
		return nil, utils.ErrBufferNotInitialized
	}

	return &ResizeInput{
		Input:        srcImg,
		TargetWidth:  c.targetWidth,
		TargetHeight: c.targetHeight,
	}, nil
}
