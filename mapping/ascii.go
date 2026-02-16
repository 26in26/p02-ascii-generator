package mapping

import (
	"fmt"
	"math"
	"strings"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/stages/edge"
)

const DENSITY = " .-=+*x#$&X@"

func pixelToASCII(gray byte, invert bool) string {
	index := int((float64(gray) / 255) * float64(len(DENSITY)-1))
	if invert {
		index = (len(DENSITY) - 1) - index
	}
	return string(DENSITY[index])
}

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

func PrintAsASCII(img *image.Buffer, gradient edge.Gradient) {
	bpp := img.Channels
	var asciiArt strings.Builder
	asciiArt.Grow((img.Width + 1) * img.Height)

	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			gradientx := y*img.Width + x
			gx := gradient[gradientx].X
			gy := gradient[gradientx].Y
			if gx*gx+gy*gy > 16 {
				angle := math.Atan2(float64(gy), float64(gx)) * 180.0 / math.Pi
				asciiArt.WriteByte(getAngleChar(angle))
			} else {
				imgx := (y*img.Width + x) * bpp
				asciiArt.WriteString(pixelToASCII(img.Data[imgx], false))
			}
		}
		asciiArt.WriteByte('\n')
	}

	fmt.Println(asciiArt.String())
}
