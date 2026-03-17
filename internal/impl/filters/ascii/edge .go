package ascii

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/pipeline/flow"
	"github.com/26in26/p02-ascii-generator/utils"
)

type EdgeFilter struct {
	pipeline.Filter[flow.Pair[*image.AsciiBuffer, utils.Gradient], *image.AsciiBuffer]
	squareThreshold int
}

func NewEdgeFilter(opts ...edgeFilterOptFunc) pipeline.Filter[flow.Pair[*image.AsciiBuffer, utils.Gradient], *image.AsciiBuffer] {
	o := defaultEdgeOpts()

	return &EdgeFilter{
		Filter:          pipeline.NewBaseFilter[flow.Pair[*image.AsciiBuffer, utils.Gradient], *image.AsciiBuffer]("Edge"),
		squareThreshold: o.threshold * o.threshold,
	}
}

func (f *EdgeFilter) Apply(ctx context.Context, input flow.Pair[*image.AsciiBuffer, utils.Gradient]) (*image.AsciiBuffer, error) {
	asciiArt := input.A
	gradient := input.B

	gradientIndex := 0
	ArtIndex := 0

	for y := 0; y < asciiArt.Height; y++ {
		for x := 0; x < asciiArt.Width; x++ {
			gradX := gradient[gradientIndex].X
			gradY := gradient[gradientIndex].Y

			if gradX*gradX+gradY*gradY > f.squareThreshold {
				asciiArt.Data[ArtIndex] = getAngleChar(gradX, gradY)
			}

			gradientIndex++
			ArtIndex += 4
		}
	}

	return asciiArt, nil
}

func getAngleChar(gx, gy int) byte {
	// Fast angle approximation using integer math to avoid expensive Atan2.
	// We classify the gradient into one of four directions.
	abs_gx := gx
	if abs_gx < 0 {
		abs_gx = -abs_gx
	}
	abs_gy := gy
	if abs_gy < 0 {
		abs_gy = -abs_gy
	}

	// If gradient is mostly horizontal -> vertical edge
	if abs_gx > abs_gy<<1 { // abs_gx > 2 * abs_gy
		return '|'
	}
	// If gradient is mostly vertical -> horizontal edge
	if abs_gy > abs_gx<<1 { // abs_gy > 2 * abs_gx
		return '_'
	}
	// Otherwise, it's a diagonal edge
	if (gx > 0) == (gy > 0) { // gx and gy have the same sign
		return '/'
	}
	return '\\' // gx and gy have different signs
}

type edgeFilterOpts struct {
	threshold int
}

type edgeFilterOptFunc func(*edgeFilterOpts)

func defaultEdgeOpts() edgeFilterOpts {
	return edgeFilterOpts{
		threshold: 20,
	}
}

func WithEdgeThreshold(threshold int) edgeFilterOptFunc {
	return func(o *edgeFilterOpts) {
		o.threshold = threshold
	}
}
