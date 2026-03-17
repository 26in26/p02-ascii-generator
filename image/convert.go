package image

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func (b *RGBBuffer) ToGray(grayBuffer *GrayBuffer) (*GrayBuffer, error) {
	dst := grayBuffer

	bpp := 3
	srcData := b.Data
	dstData := dst.Data

	dstIndex := 0
	end := len(srcData) - bpp

	for i := 0; i <= end; i += bpp {
		r := int(srcData[i])
		g := int(srcData[i+1])
		b := int(srcData[i+2])

		// Convert to grayscale using integer approximation (0.21 R + 0.72 G + 0.07 B)
		// Scaled by 256: 54 R + 184 G + 18 B. Sum = 256.
		gray := (r*54 + g*184 + b*18) >> 8
		dstData[dstIndex] = uint8(gray)

		dstIndex++
	}

	return dst, nil
}

func (b *AsciiBuffer) ToImage() image.Image {

	const cellW = 7
	const cellH = 13

	imgWidth := b.buffer.Width * cellW
	imgHeight := b.buffer.Height * cellH

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: basicfont.Face7x13,
	}

	for y := 0; y < b.buffer.Height; y++ {
		for x := 0; x < b.buffer.Width; x++ {

			i := y*b.buffer.Stride + x*4

			ch := rune(b.buffer.Data[i])
			r := b.buffer.Data[i+1]
			g := b.buffer.Data[i+2]
			b := b.buffer.Data[i+3]

			drawer.Src = image.NewUniform(color.RGBA{r, g, b, 255})

			px := x * cellW
			py := y * cellH

			// glyphs are drawn from the basline up (not top-left corner)
			drawer.Dot = fixed.P(px, py+13)

			drawer.DrawString(string(ch))
		}
	}

	return img
}
