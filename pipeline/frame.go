package pipeline

import (
	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/utils"
)

type FrameContext struct {
	SourceImage  *image.RGBBuffer
	WorkingImage *image.RGBBuffer
	GrayImage    *image.GrayBuffer
	GradientMap  utils.Gradient
	AsciiArt     *image.AsciiBuffer
}

func NewFrameContext(src *image.RGBBuffer) *FrameContext {
	return &FrameContext{
		SourceImage: src,
	}
}
