package ascii

import (
	"fmt"
	"math"
	"strings"

	"github.com/26in26/p02-ascii-generator/pipeline"
)

const DENSITY = " .-=+*x#$&X@"
const RESET = "\x1b[0m"

type ColorMode uint8

const (
	FullColor ColorMode = iota
	RetroColor
	Monochrome
)

func getAngleChar(angle float64) byte {
	if (22.5 <= angle && angle <= 67.5) || (-157.5 <= angle && angle <= -112.5) {
		return '\\'
	} else if (67.5 <= angle && angle <= 112.5) || (-112.5 <= angle && angle <= -67.5) {
		return '_'
	} else if (112.5 <= angle && angle <= 157.5) || (-67.5 <= angle && angle <= -22.5) {
		return '/'
	} else {
		return '|'
	}
}

func pixelToASCII(gray byte, invert bool) byte {
	index := int((float64(gray) / 255) * float64(len(DENSITY)-1))
	if invert {
		index = (len(DENSITY) - 1) - index
	}
	return DENSITY[index]
}

type AsciiStage struct {
	invert              bool
	squareEdgeThreshold int
}

func NewAsciiStage(invert bool, edgeThreshold int) *AsciiStage {
	return &AsciiStage{
		invert:              invert,
		squareEdgeThreshold: edgeThreshold * edgeThreshold,
	}
}

func (s *AsciiStage) Process(ctx *pipeline.FrameContext) {
	workingImg := ctx.WorkingImage
	if workingImg == nil {
		panic("ascii stage must receive non-nil image buffer")
	}

	grayImg := ctx.GrayImage
	gradient := ctx.GradientMap

	bpp := workingImg.Channels
	var asciiArt strings.Builder
	asciiArt.Grow((workingImg.Width + 1) * workingImg.Height)

	var r, g, b byte = 0, 0, 0

	for y := 0; y < grayImg.Height; y++ {
		for x := 0; x < grayImg.Width; x++ {
			var char byte

			index := y*grayImg.Width + x

			gx := gradient[index].X
			gy := gradient[index].Y
			if gx*gx+gy*gy > s.squareEdgeThreshold {
				angle := math.Atan2(float64(gy), float64(gx)) * 180.0 / math.Pi
				char = getAngleChar(angle)
			} else {
				char = pixelToASCII(grayImg.Data[index], s.invert)
			}

			if r != workingImg.Data[index*bpp+0] || g != workingImg.Data[index*bpp+1] || b != workingImg.Data[index*bpp+2] {
				r = workingImg.Data[index*bpp+0]
				g = workingImg.Data[index*bpp+1]
				b = workingImg.Data[index*bpp+2]

				asciiArt.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b))
			}

			asciiArt.WriteByte(char)

		}

		asciiArt.WriteByte('\n')
	}
	asciiArt.WriteString(RESET)

	ctx.ASCIIOutput = asciiArt.String()
	println(ctx.ASCIIOutput)
}
