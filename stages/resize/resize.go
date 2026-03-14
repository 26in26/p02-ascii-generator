package resize

import (
	"context"
	"fmt"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/utils"
)

const CHAR_ASPECT_RATIO = 2

type ResizeStage struct {
	pipeline.BaseStage[*image.RGBBuffer, *image.RGBBuffer]
	w int
	h int
}

func NewResizeStage(opts ...optFunc) (pipeline.Stage[*image.RGBBuffer, *image.RGBBuffer], error) {
	o := defaultOpts()

	for _, opt := range opts {
		opt(&o)
	}

	if o.width <= 0 || o.height <= 0 {
		return nil, fmt.Errorf("resize stage: %w", utils.ErrInvalidDimensions)
	}

	return &ResizeStage{
		BaseStage: *pipeline.NewBaseStage[*image.RGBBuffer, *image.RGBBuffer]("Resize"),
		w:         o.width,
		h:         o.height,
	}, nil
}

func (s *ResizeStage) Kernal(ctx context.Context, input *image.RGBBuffer) (*image.RGBBuffer, error) {
	src := input

	dst, err := image.NewRGBBuffer(s.w, s.h)

	if err != nil {
		return nil, err
	}

	bpp := 3

	xRatio := float64(src.Width) / float64(dst.Width)
	yRatio := float64(src.Height) / float64(dst.Height)

	srcXOffsets := make([]int, dst.Width)
	for x := 0; x < dst.Width; x++ {
		srcX := int(float64(x) * xRatio)
		if srcX >= src.Width {
			srcX = src.Width - 1
		}
		srcXOffsets[x] = srcX * bpp
	}

	srcStride := src.Stride
	dstStride := dst.Stride
	srcData := src.Data
	dstData := dst.Data

	for y := 0; y < dst.Height; y++ {
		srcY := int(float64(y) * yRatio)
		if srcY >= src.Height {
			srcY = src.Height - 1
		}

		srcRowOffset := srcY * srcStride
		dstOffset := y * dstStride

		for x := 0; x < dst.Width; x++ {
			srcOffset := srcRowOffset + srcXOffsets[x]

			copy(
				dstData[dstOffset:dstOffset+bpp],
				srcData[srcOffset:srcOffset+bpp],
			)
			dstOffset += bpp
		}
	}

	return dst, nil
}

func (s *ResizeStage) Release(input *image.RGBBuffer) {
	if input == nil {
		return
	}
}
