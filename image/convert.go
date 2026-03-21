package image

import (
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"strings"

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

func (a *AsciiBuffer) ToImage(useColor bool) image.Image {
	var img image.Image
	if useColor {
		img = a.ToImageWithFullColor()
	} else {
		img = a.ToImageWithSingleColor(255, 255, 255)
	}

	return img
}

func (a *AsciiBuffer) ToImageWithFullColor() image.Image {

	const cellW = 7
	const cellH = 13

	imgWidth := a.buffer.Width * cellW
	imgHeight := a.buffer.Height * cellH

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	draw.Draw(
		img,
		img.Bounds(),
		&image.Uniform{color.Black},
		image.Point{},
		draw.Src,
	)

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: basicfont.Face7x13,
	}

	index := 0

	for y := 0; y < a.buffer.Height; y++ {
		for x := 0; x < a.buffer.Width; x++ {

			ch := rune(a.buffer.Data[index])
			r := a.buffer.Data[index+1]
			g := a.buffer.Data[index+2]
			b := a.buffer.Data[index+3]

			drawer.Src = image.NewUniform(color.RGBA{r, g, b, 255})

			px := x * cellW
			py := y * cellH

			// glyphs are drawn from the basline up (not top-left corner)
			drawer.Dot = fixed.P(px, py+13)

			drawer.DrawString(string(ch))

			index += 4
		}
	}

	return img
}

func (a *AsciiBuffer) ToImageWithSingleColor(r, g, b byte) image.Image {
	const cellW = 7
	const cellH = 13

	imgWidth := a.buffer.Width * cellW
	imgHeight := a.buffer.Height * cellH

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(
		img,
		img.Bounds(),
		&image.Uniform{color.Black},
		image.Point{},
		draw.Src,
	)

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.NRGBA{r, g, b, 255}),
		Face: basicfont.Face7x13,
	}

	index := 0

	for y := 0; y < a.buffer.Height; y++ {
		for x := 0; x < a.buffer.Width; x++ {

			ch := rune(a.buffer.Data[index])

			px := x * cellW
			py := y * cellH

			// glyphs are drawn from the basline up (not top-left corner)
			drawer.Dot = fixed.P(px, py+13)

			drawer.DrawString(string(ch))

			index += 4
		}
	}

	return img
}

const RESET = "\x1b[0m"

func (a *AsciiBuffer) ToString(str *strings.Builder, useColor bool) {
	if useColor {
		a.ToStringWithColor(str)
	} else {
		a.ToStringWithoutColor(str)
	}
}

func (a *AsciiBuffer) ToStringWithColor(str *strings.Builder) {
	i := 0
	var curR, curG, curB byte

	for y := 0; y < a.buffer.Height; y++ {
		for x := 0; x < a.buffer.Width; x++ {
			char, r, g, b := a.Data[i], a.Data[i+1], a.Data[i+2], a.Data[i+3]
			printColor(str, curR, curG, curB, r, g, b)
			curR, curG, curB = r, g, b
			str.WriteByte(char)
			i += 4
		}
		str.WriteByte('\n')
	}
	str.WriteString(RESET)
}

func (a *AsciiBuffer) ToStringWithoutColor(str *strings.Builder) {
	i := 0

	for y := 0; y < a.buffer.Height; y++ {
		for x := 0; x < a.buffer.Width; x++ {
			char := a.Data[i]
			str.WriteByte(char)
			i += 4
		}
		str.WriteByte('\n')
	}
	str.WriteString(RESET)
}

// renderFullColor writes a 24-bit ANSI color escape code if the color has changed.
// curR, curG, curB := data[rgbIndex], data[rgbIndex+1], data[rgbIndex+2]

func printColor(builder *strings.Builder, curR, curG, curB, r, g, b byte) {

	if r != curR || g != curG || b != curB {
		builder.WriteString("\x1b[38;2;")
		builder.WriteString(strconv.Itoa(int(r)))
		builder.WriteByte(';')
		builder.WriteString(strconv.Itoa(int(g)))
		builder.WriteByte(';')
		builder.WriteString(strconv.Itoa(int(b)))
		builder.WriteByte('m')
	}
}
