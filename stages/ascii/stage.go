package ascii

import (
	"context"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/utils"
)

type AsciiStage struct {
	invert              bool
	squareEdgeThreshold int
}

type AsciiInput struct {
	workingImg *image.RGBBuffer
	grayImg    *image.GrayBuffer
	gradient   utils.Gradient
}

func NewAsciiStage(opts ...optFunc) pipeline.Stage {
	o := defaultOpts()

	for _, opt := range opts {
		opt(&o)
	}

	// s := &AsciiStage{
	// 	invert:              o.invert,
	// 	squareEdgeThreshold: o.edgeThreshold * o.edgeThreshold,
	// }

	return pipeline.NewBaseStage(
		"ascii", []pipeline.DataType{pipeline.DataResized, pipeline.DataAscii, pipeline.DataGray}, pipeline.DataAscii, NewAsciiConnector(), Ascii,
	)
}

func Ascii(ctx context.Context, input *AsciiInput) (*image.AsciiBuffer, error) {
	workingData := input.workingImg.Data
	grayData := input.grayImg.Data
	gradient := input.gradient

	bpp := 3
	asciiArt, err := image.NewAsciiBuffer(input.workingImg.Width, input.workingImg.Height)
	if err != nil {
		return nil, err
	}

	var r, g, b byte = 0, 0, 0
	index := 0

	for y := 0; y < asciiArt.Height; y++ {
		for x := 0; x < asciiArt.Width; x++ {
			var char byte
			if gradient[index].X*gradient[index].X+gradient[index].Y*gradient[index].Y > 400 {
				char = getAngleChar(gradient[index].X, gradient[index].Y)
			} else {
				char = pixelToASCII(grayData[index], false)
			}

			r, g, b = workingData[index*bpp], workingData[index*bpp+1], workingData[index*bpp+2]
			asciiArt.Data[index*4] = char
			asciiArt.Data[index*4+1] = r
			asciiArt.Data[index*4+2] = g
			asciiArt.Data[index*4+3] = b

			index++
		}
	}

	return asciiArt, nil
}
