package image

import "github.com/26in26/p02-ascii-generator/utils"

type Buffer struct {
	Width  int
	Height int
	stride int

	Format   Format
	Channels int
	Data     []byte
}

func NewBuffer(width, height int, format Format) (*Buffer, error) {
	if width < 0 || height < 0 {
		return nil, utils.ErrInvalidDimensions
	}

	bpp := format.BytesPerPixel()
	stride := width * bpp

	return &Buffer{
		Width:    width,
		Height:   height,
		stride:   stride,
		Format:   format,
		Channels: bpp,
		Data:     make([]byte, stride*height),
	}, nil

}

func (b *Buffer) Stride() int {
	return b.stride
}
