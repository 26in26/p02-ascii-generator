package ascii

import (
	"fmt"
	"strings"

	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/utils"
)

type charSelector func(index int, gray byte, gradient utils.Vector) byte

type colorRenderer func(builder *strings.Builder, r, g, b byte, data []byte, rgbIndex int) (byte, byte, byte)

type AsciiStage struct {
	invert              bool
	squareEdgeThreshold int

	selector charSelector
	renderer colorRenderer
}

func NewAsciiStage(opts ...optFunc) *AsciiStage {
	o := defaultOpts()

	for _, opt := range opts {
		opt(&o)
	}

	s := &AsciiStage{
		invert:              o.invert,
		squareEdgeThreshold: o.edgeThreshold * o.edgeThreshold,
	}

	if o.useEdge {
		s.selector = s.selectEdgeChar
	} else {
		s.selector = s.selectDensityChar
	}

	// Set up the color renderer based on options.
	switch o.colorMode {
	case FullColor:
		s.renderer = renderFullColor
	case Monochrome:
		s.renderer = renderMonochrome
	default:
		s.renderer = renderMonochrome
	}

	return s
}

func (s *AsciiStage) selectEdgeChar(index int, gray byte, gradient utils.Vector) byte {
	if gradient.X*gradient.X+gradient.Y*gradient.Y > s.squareEdgeThreshold {
		return getAngleChar(gradient.X, gradient.Y)
	}
	return pixelToASCII(gray, s.invert)
}

func (s *AsciiStage) selectDensityChar(index int, gray byte, gradient utils.Vector) byte {
	return pixelToASCII(gray, s.invert)
}

func (s *AsciiStage) Process(ctx *pipeline.FrameContext) error {
	workingImg := ctx.WorkingImage
	if workingImg == nil {
		return fmt.Errorf("ascii stage: %w", utils.ErrBufferNotInitialized)
	}

	grayImg := ctx.GrayImage
	gradient := ctx.GradientMap

	bpp := workingImg.Channels
	var asciiArt strings.Builder
	asciiArt.Grow((workingImg.Width + 1) * workingImg.Height)

	grayData := grayImg.Data
	workingData := workingImg.Data

	var r, g, b byte = 0, 0, 0
	index := 0

	for y := 0; y < grayImg.Height; y++ {
		for x := 0; x < grayImg.Width; x++ {
			char := s.selector(index, grayData[index], gradient[index])

			r, g, b = s.renderer(&asciiArt, r, g, b, workingData, index*bpp)

			asciiArt.WriteByte(char)
			index++
		}

		asciiArt.WriteByte('\n')
	}
	asciiArt.WriteString(RESET)

	ctx.ASCIIOutput = asciiArt.String()
	println(ctx.ASCIIOutput)

	return nil
}
