package grayscale

import (
	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/utils/pools"
)

type size struct {
	w, h int
}

type grayImagePool struct {
	*pools.KeyedPool[size, *image.GrayBuffer]
}

func NewGrayImagePool() *grayImagePool {
	return &grayImagePool{
		KeyedPool: pools.NewKeyedPool(func(s size) (*image.GrayBuffer, error) {
			return image.NewGrayBuffer(s.w, s.h)
		}),
	}
}
