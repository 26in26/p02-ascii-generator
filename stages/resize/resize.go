package resize

import "github.com/26in26/p02-ascii-generator/image"

const CHAR_ASPECT_RATIO = 2

type ResizeStage struct {
	TargetWidth  int
	TargetHeight int
}

func NewResizeStage(w, h int) *ResizeStage {
	if w <= 0 || h <= 0 {
		panic("resize stage width and height must be > 0")
	}

	return &ResizeStage{
		TargetWidth:  w,
		TargetHeight: h,
	}
}

func (s *ResizeStage) PreserveAspectRatio(w, h int, saveWidth, saveHeight bool,
) *ResizeStage {

	aspectRatio := float64(w) / float64(h) * CHAR_ASPECT_RATIO

	if saveWidth {
		// adjust height according to aspect ratio
		s.TargetHeight = int(float64(s.TargetWidth) / aspectRatio)
		if s.TargetHeight <= 0 {
			s.TargetHeight = 1
		}
	} else if saveHeight {
		// adjust width according to aspect ratio
		s.TargetWidth = int(float64(s.TargetHeight) * aspectRatio)
		if s.TargetWidth <= 0 {
			s.TargetWidth = 1
		}
	}
	return s
}

func (s *ResizeStage) Process(src *image.Buffer) *image.Buffer {
	if src == nil {
		panic("resize stage must receive non-nil image buffer")
	}

	dst := image.NewBuffer(s.TargetWidth, s.TargetHeight, src.Format)

	bpp := src.Channels

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

	srcStride := src.Stride()
	dstStride := dst.Stride()
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

	return dst
}
