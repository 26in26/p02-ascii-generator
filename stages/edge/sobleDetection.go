package edge

import (
	"github.com/26in26/p02-ascii-generator/image"
)

type Vector struct {
	X int
	Y int
}

type Gradient = []Vector

type sobelEdgeDetectionStage struct {
	Gradient Gradient
}

func NewSobelEdgeDetectionStage() *sobelEdgeDetectionStage {
	return &sobelEdgeDetectionStage{}
}

func (s *sobelEdgeDetectionStage) Process(src *image.Buffer) *image.Buffer {
	if src == nil {
		panic("sobel edge detection stage must receive non-nil image buffer")
	}

	// For performance, we first convert the image to grayscale.
	if src.Format != image.FormatGray {
		src = src.ToGray()
	}

	srcData := src.Data
	width := src.Width
	height := src.Height
	s.Gradient = make(Gradient, len(srcData))
	gradData := s.Gradient

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

			// Apply Sobel kernels directly. Using bit-shifts (<< 1) for
			// multiplication by 2 is a common performance pattern.
			// scaling down the vectors by 16 ( >> 4)
			gx := ((int(p2) + (int(p5) << 1) + int(p8)) - (int(p0) + (int(p3) << 1) + int(p6))) >> 4
			gy := ((int(p6) + (int(p7) << 1) + int(p8)) - (int(p0) + (int(p1) << 1) + int(p2))) >> 4
			gradData[i] = Vector{gx, gy}

			i++
		}

		// Skip the right border of the current row and the left border of the next row
		i += 2
	}

	return src
}
