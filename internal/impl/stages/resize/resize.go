package resize

import (
	"context"
	"fmt"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/internal/errs"
	"github.com/26in26/p02-ascii-generator/pipeline"
)

const CHAR_ASPECT_RATIO = 2

type resizeStage struct {
	pipeline.BaseStage[*image.RGBBuffer, *image.RGBBuffer]
	imagePool *RGBImagePool
	xOffsets  *xOffsetsPool
	w         int
	h         int
}

func NewResizeStage(w, h int) (pipeline.Stage[*image.RGBBuffer, *image.RGBBuffer], error) {
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("resize stage: %w", errs.InvalidDimensions)
	}

	return &resizeStage{
		BaseStage: *pipeline.NewBaseStage[*image.RGBBuffer, *image.RGBBuffer]("Resize"),
		imagePool: NewRGBImagePool(w, h),
		xOffsets:  newXOffsetsPool(w),
		w:         w,
		h:         h,
	}, nil
}

func (s *resizeStage) Kernal(ctx context.Context, input *image.RGBBuffer) (*image.RGBBuffer, error) {
	src := input

	dst, err := s.imagePool.Get()

	if err != nil {
		return nil, err
	}

	bpp := 3

	xRatio := float64(src.Width) / float64(dst.Width)
	yRatio := float64(src.Height) / float64(dst.Height)

	srcXOffsets := s.xOffsets.Get()
	defer s.xOffsets.Put(srcXOffsets)

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

func (s *resizeStage) Release(input *image.RGBBuffer) {
	if input == nil {
		return
	}
}
