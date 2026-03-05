package edge

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/utils"

	"github.com/26in26/p02-ascii-generator/pipeline"
)

type EdgeConnector struct{}

func NewEdgeConnector() pipeline.StageConnector[*image.GrayBuffer, utils.Gradient] {
	return &EdgeConnector{}
}

func (c *EdgeConnector) Set(ctx context.Context, f *pipeline.Frame, val utils.Gradient) error {
	f.GradientMap.Set(val)
	return nil
}

func (c *EdgeConnector) Get(ctx context.Context, f *pipeline.Frame) (*image.GrayBuffer, error) {
	srcImg, err := f.GrayImage.Get(ctx)

	if err != nil {
		return nil, err
	}

	if srcImg == nil {
		return nil, utils.ErrBufferNotInitialized
	}

	return srcImg, nil
}
