package image

import (
	"image"
	"image/color"
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

const RESET = "\x1b[0m"

func (b *AsciiBuffer) ToString(str *strings.Builder, useColor bool) {
	if useColor {
		b.ToStringWithColor(str)
	} else {
		b.ToStringWithoutColor(str)
	}
}

func (b *AsciiBuffer) ToStringWithColor(str *strings.Builder) {
	i := 0
	var curR, curG, curB byte

	for y := 0; y < b.buffer.Height; y++ {
		for x := 0; x < b.buffer.Width; x++ {
			char, r, g, b := b.Data[i], b.Data[i+1], b.Data[i+2], b.Data[i+3]
			printColor(str, curR, curG, curB, r, g, b)
			curR, curG, curB = r, g, b
			str.WriteByte(char)
			i += 4
		}
		str.WriteByte('\n')
	}
	str.WriteString(RESET)
}

func (b *AsciiBuffer) ToStringWithoutColor(str *strings.Builder) {
	i := 0

	for y := 0; y < b.buffer.Height; y++ {
		for x := 0; x < b.buffer.Width; x++ {
			char := b.Data[i]
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
