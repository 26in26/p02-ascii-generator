package imageio

import (
	"image"
	"image/color"

	internalImage "github.com/26in26/p02-ascii-generator/image"
)

// Todo: support formats, currently only RGB & RGBA
func ConvertToBuffer(img image.Image, f internalImage.Format) (*internalImage.Buffer, error) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	buf, err := internalImage.NewBuffer(w, h, f)

	if err != nil {
		return nil, err
	}

	bpp := buf.Channels

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			idx := ((y-bounds.Min.Y)*w + (x - bounds.Min.X)) * bpp
			buf.Data[idx+0] = c.R
			buf.Data[idx+1] = c.G
			buf.Data[idx+2] = c.B
			if bpp == 4 {
				buf.Data[idx+3] = c.A
			}
		}
	}

	return buf, nil
}
