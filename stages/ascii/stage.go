package ascii

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/pipeline/flow"
	"github.com/26in26/p02-ascii-generator/utils"
)

type AsciiStage struct {
	*pipeline.BaseStage[*image.GrayBuffer, *image.AsciiBuffer]
	invert bool
}

func NewAsciiStage(opts ...stageOptFunc) pipeline.Stage[*image.GrayBuffer, *image.AsciiBuffer] {
	o := defaultOpts()

	for _, opt := range opts {
		opt(&o)
	}

	return &AsciiStage{
		BaseStage: pipeline.NewBaseStage[*image.GrayBuffer, *image.AsciiBuffer]("Ascii"),
		invert:    o.invert,
	}
}

func (s *AsciiStage) Kernal(ctx context.Context, input *image.GrayBuffer) (*image.AsciiBuffer, error) {
	asciiArt, err := image.NewAsciiBuffer(input.Width, input.Height)
	if err != nil {
		return nil, err
	}

	grayIndex := 0
	ArtIndex := 0

	for y := 0; y < asciiArt.Height; y++ {
		for x := 0; x < asciiArt.Width; x++ {
			asciiArt.Data[ArtIndex] = pixelToASCII(input.Data[grayIndex], s.invert)

			grayIndex++
			ArtIndex += 4
		}
	}

	return asciiArt, nil
}

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

type ColorFilter struct {
	pipeline.Filter[flow.Pair[*image.AsciiBuffer, *image.RGBBuffer], *image.AsciiBuffer]
	colorMode ColorMode
}

func NewColorFilter(opts ...colorFilterOptFunc) pipeline.Filter[flow.Pair[*image.AsciiBuffer, *image.RGBBuffer], *image.AsciiBuffer] {
	o := defaultColorOpts()

	for _, opt := range opts {
		opt(&o)
	}

	return &ColorFilter{
		Filter:    pipeline.NewBaseFilter[flow.Pair[*image.AsciiBuffer, *image.RGBBuffer], *image.AsciiBuffer]("Color"),
		colorMode: o.colorMode,
	}
}

func (f *ColorFilter) Apply(ctx context.Context, input flow.Pair[*image.AsciiBuffer, *image.RGBBuffer]) (*image.AsciiBuffer, error) {
	asciiArt := input.A
	RGBData := input.B.Data

	artIndex := 0
	RGBIndex := 0

	for y := 0; y < asciiArt.Height; y++ {
		for x := 0; x < asciiArt.Width; x++ {
			r, g, b := RGBData[RGBIndex], RGBData[RGBIndex+1], RGBData[RGBIndex+2]
			asciiArt.Data[artIndex+1] = r
			asciiArt.Data[artIndex+2] = g
			asciiArt.Data[artIndex+3] = b

			artIndex += 4
			RGBIndex += 3

		}
	}

	return asciiArt, nil
}
