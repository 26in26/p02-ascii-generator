package ascii

import (
	"strconv"
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
		return '\\'
	}
	return '/' // gx and gy have different signs
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

	grayData := grayImg.Data
	workingData := workingImg.Data

	var r, g, b byte = 0, 0, 0
	index := 0

	for y := 0; y < grayImg.Height; y++ {
		for x := 0; x < grayImg.Width; x++ {
			var char byte

			gx := gradient[index].X
			gy := gradient[index].Y
			if gx*gx+gy*gy > s.squareEdgeThreshold {
				char = getAngleChar(gx, gy)
			} else {
				char = pixelToASCII(grayData[index], s.invert)
			}

			rgbIndex := index * bpp
			curR, curG, curB := workingData[rgbIndex], workingData[rgbIndex+1], workingData[rgbIndex+2]

			if r != curR || g != curG || b != curB {
				r, g, b = curR, curG, curB

				asciiArt.WriteString("\x1b[38;2;")
				asciiArt.WriteString(strconv.Itoa(int(r)))
				asciiArt.WriteByte(';')
				asciiArt.WriteString(strconv.Itoa(int(g)))
				asciiArt.WriteByte(';')
				asciiArt.WriteString(strconv.Itoa(int(b)))
				asciiArt.WriteByte('m')
			}

			asciiArt.WriteByte(char)
			index++
		}

		asciiArt.WriteByte('\n')
	}
	asciiArt.WriteString(RESET)

	ctx.ASCIIOutput = asciiArt.String()
	println(ctx.ASCIIOutput)
}
