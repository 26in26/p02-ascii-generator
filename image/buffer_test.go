package image_test

import (
	"errors"
	"testing"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/internal/errs"
)

func TestNewBufferDimensionsAndStride(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		stride int
	}{{
		name:   "positive case",
		width:  10,
		height: 5,
		stride: 30,
	}, {
		name:   "zero height",
		width:  10,
		height: 0,
		stride: 30,
	}, {
		name:   "zero width",
		width:  0,
		height: 15,
		stride: 0,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf, _ := image.NewRGBBuffer(tc.width, tc.height)

			if buf.Width != tc.width {
				t.Fatalf("width = %d, want %d", buf.Width, tc.width)
			}
			if buf.Height != tc.height {
				t.Fatalf("height = %d, want %d", buf.Height, tc.height)
			}
			if buf.Stride != tc.stride {
				t.Fatalf("stride = %d, want %d", buf.Stride, tc.stride)
			}
		})

	}

}

func TestInvalidDimensions(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{{
		name:   "negative height",
		width:  10,
		height: -5,
	}, {
		name:   "negative width",
		width:  -10,
		height: 5,
	}, {
		name:   "negative width & height",
		width:  -10,
		height: -5,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			_, err := image.NewRGBBuffer(tc.width, tc.height)
			if err == nil {
				t.Fatal("expected error, got nil")
			} else if !errors.Is(err, errs.InvalidDimensions) {
				t.Fatalf("Expexeted ErrInvalidDimensions, got %v", err)
			}
		})
	}
}

func TestNewBufferDataSize(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		height     int
		targetSize int
	}{{
		name:       "positive case",
		width:      10,
		height:     5,
		targetSize: 150,
	}, {
		name:       "zero width",
		width:      0,
		height:     5,
		targetSize: 0,
	}, {
		name:       "zero height",
		width:      0,
		height:     15,
		targetSize: 0,
	}, {
		name:       "zero width & height",
		width:      0,
		height:     0,
		targetSize: 0,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf, _ := image.NewRGBBuffer(tc.width, tc.height)

			if len(buf.Data) != tc.targetSize {
				t.Fatalf("data size = %d, want %d", len(buf.Data), tc.targetSize)
			}
		})
	}

}
