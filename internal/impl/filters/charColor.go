package filters

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/pipeline/flow"
)

type ColorMode uint8

const (
	FullColor ColorMode = iota
	RetroColor
	Monochrome
)

type ColorFilter struct {
	pipeline.Filter[flow.Pair[*image.AsciiBuffer, *image.RGBBuffer], *image.AsciiBuffer]
	colorMode ColorMode
}

func NewColorFilter(colorMode ColorMode) pipeline.Filter[flow.Pair[*image.AsciiBuffer, *image.RGBBuffer], *image.AsciiBuffer] {
	return &ColorFilter{
		Filter:    pipeline.NewBaseFilter[flow.Pair[*image.AsciiBuffer, *image.RGBBuffer], *image.AsciiBuffer]("Color"),
		colorMode: colorMode,
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
