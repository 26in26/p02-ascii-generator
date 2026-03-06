package resize

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/utils"
)

type ResizeConnector struct{}

func (c *ResizeConnector) Set(ctx context.Context, f *pipeline.Frame, val *image.RGBBuffer) error {
	f.ResizedImage.Set(val)
	return nil
}

func (c *ResizeConnector) Get(ctx context.Context, f *pipeline.Frame) (*image.RGBBuffer, error) {
	srcImg, err := f.Raw.Get(ctx)

	if err != nil {
		return nil, err
	}

	if srcImg == nil {
		return nil, utils.ErrBufferNotInitialized
	}

	return srcImg, nil
}
