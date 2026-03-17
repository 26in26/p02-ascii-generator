package grayscale

import (
	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/utils"
)

type size struct {
	w, h int
}

type grayImagePool struct {
	*utils.Pool[size, *image.GrayBuffer]
}

func NewGrayImagePool() *grayImagePool {
	return &grayImagePool{
		Pool: utils.NewPool(func(s size) (*image.GrayBuffer, error) {
			return image.NewGrayBuffer(s.w, s.h)
		}),
	}
}
