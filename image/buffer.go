package image

import (
	"github.com/26in26/p02-ascii-generator/internal/errs"
)

type buffer struct {
	Width  int
	Height int
	Stride int
	Data   []byte
}

// RGBBuffer represents a 3-channel RGB image.
// It provides 3 color chanels (bytes) per pixel.
type RGBBuffer struct {
	buffer
}

func NewRGBBuffer(width, height int) (*RGBBuffer, error) {
	if width < 0 || height < 0 {
		return nil, errs.InvalidDimensions
	}
	stride := width * 3
	return &RGBBuffer{
		buffer: buffer{
			Width:  width,
			Height: height,
			Stride: stride,
			Data:   make([]byte, stride*height),
		},
	}, nil
}

// GrayBuffer represents a 1-channel Grayscale image.
// It provides 1 gray channel (byte) per pixel.
type GrayBuffer struct {
	buffer
}

func NewGrayBuffer(width, height int) (*GrayBuffer, error) {
	if width < 0 || height < 0 {
		return nil, errs.InvalidDimensions
	}
	return &GrayBuffer{
		buffer: buffer{
			Width:  width,
			Height: height,
			Stride: width,
			Data:   make([]byte, width*height),
		},
	}, nil
}

// AsciiBuffer represents an image made from colored text.
// It provides 1 char and 3 color chanels (bytes) per pixel.
type AsciiBuffer struct {
	buffer
}

func NewAsciiBuffer(width, height int) (*AsciiBuffer, error) {
	if width < 0 || height < 0 {
		return nil, errs.InvalidDimensions
	}

	stride := width * 4

	return &AsciiBuffer{
		buffer: buffer{
			Width:  width,
			Height: height,
			Stride: stride,
			Data:   make([]byte, stride*height),
		},
	}, nil
}
