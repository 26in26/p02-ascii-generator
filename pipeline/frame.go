package pipeline

import (
	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/utils"
)

type FrameContext struct {
	SourceImage  *image.Buffer
	WorkingImage *image.Buffer
	GrayImage    *image.Buffer
	GradientMap  utils.Gradient
	ASCIIOutput  string
}

func NewFrameContext(src *image.Buffer) *FrameContext {
	return &FrameContext{
		SourceImage: src,
	}
}
