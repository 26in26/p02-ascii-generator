package edge

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/utils"
)

type SobelEdgeDetectionStage struct {
	pipeline.BaseStage[*image.GrayBuffer, utils.Gradient]
	gradPool *gradientPool
}

func NewSobelEdgeDetectionStage() pipeline.Stage[*image.GrayBuffer, utils.Gradient] {
	return &SobelEdgeDetectionStage{
		BaseStage: *pipeline.NewBaseStage[*image.GrayBuffer, utils.Gradient]("SobelEdgeDetection"),
		gradPool:  NewGrayImagePool(),
	}
}

func (s *SobelEdgeDetectionStage) Kernal(ctx context.Context, input *image.GrayBuffer) (utils.Gradient, error) {
	srcData := input.Data
	width := input.Width
	height := input.Height

	gradData, err := s.gradPool.Get(len(srcData))
	if err != nil {
		return nil, err
	}

	i := width + 1

	for y := 1; y < height-1; y++ {

		for x := 1; x < width-1; x++ {
			// access all 8 pixel neighbors without any checks. This is very fast.
			p0 := srcData[i-width-1] // top-left
			p1 := srcData[i-width]   // top
			p2 := srcData[i-width+1] // top-right
			p3 := srcData[i-1]       // left
			p5 := srcData[i+1]       // right
			p6 := srcData[i+width-1] // bottom-left
			p7 := srcData[i+width]   // bottom
			p8 := srcData[i+width+1] // bottom-right

			// Apply Sobel kernels directly. Using bit-shifts (<< 1) for preformance
			// scaling down the vectors by 16 ( >> 4)
			gx := ((int(p2) + (int(p5) << 1) + int(p8)) - (int(p0) + (int(p3) << 1) + int(p6))) >> 4
			gy := ((int(p6) + (int(p7) << 1) + int(p8)) - (int(p0) + (int(p1) << 1) + int(p2))) >> 4
			gradData[i] = utils.Vector{X: gx, Y: gy}

			i++
		}

		// Skip the right border of the current row and the left border of the next row
		i += 2
	}

	return gradData, nil
}

func (s *SobelEdgeDetectionStage) Release(grad utils.Gradient) {
	if grad == nil {
		return
	}

	s.gradPool.Put(len(grad), grad)
}
