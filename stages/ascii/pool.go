package ascii

import (
	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/utils"
)

type size struct {
	w, h int
}

type asciiArtPool struct {
	*utils.Pool[size, *image.AsciiBuffer]
}

func NewAsciiArtPool() *asciiArtPool {
	return &asciiArtPool{
		Pool: utils.NewPool(func(s size) (*image.AsciiBuffer, error) {
			return image.NewAsciiBuffer(s.w, s.h)
		}),
	}
}
