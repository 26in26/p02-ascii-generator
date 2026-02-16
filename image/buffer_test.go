package image_test

import (
	"testing"

	"github.com/26in26/p02-ascii-generator/image"
)

func TestNewBufferDimensionsAndStride(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		format image.Format
		stride int
	}{{
		name:   "positive case",
		width:  10,
		height: 5,
		format: image.FormatRGB,
		stride: 30,
	}, {
		name:   "zero height",
		width:  10,
		height: 0,
		format: image.FormatRGB,
		stride: 30,
	}, {
		name:   "zero width",
		width:  0,
		height: 15,
		format: image.FormatRGB,
		stride: 0,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := image.NewBuffer(tc.width, tc.height, tc.format)

			if buf.Width != tc.width {
				t.Fatalf("width = %d, want %d", buf.Width, tc.width)
			}
			if buf.Height != tc.height {
				t.Fatalf("height = %d, want %d", buf.Height, tc.height)
			}
			if buf.Stride() != tc.stride {
				t.Fatalf("stride = %d, want %d", buf.Stride(), tc.stride)
			}
		})

	}

}

func TestInvalidDimensionsPanics(t *testing.T) {
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
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("expected panic due to invalid dimensions, but got none")
				}
			}()

			image.NewBuffer(tc.width, tc.height, image.FormatRGB)
		})
	}
}

func TestNewBufferDataSize(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		height     int
		format     image.Format
		targetSize int
	}{{
		name:       "positive case",
		width:      10,
		height:     5,
		format:     image.FormatRGB,
		targetSize: 150,
	}, {
		name:       "zero width",
		width:      0,
		height:     5,
		format:     image.FormatRGB,
		targetSize: 0,
	}, {
		name:       "zero height",
		width:      0,
		height:     15,
		format:     image.FormatRGB,
		targetSize: 0,
	}, {
		name:       "zero width & height",
		width:      0,
		height:     0,
		format:     image.FormatRGB,
		targetSize: 0,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := image.NewBuffer(tc.width, tc.height, tc.format)

			if len(buf.Data) != tc.targetSize {
				t.Fatalf("data size = %d, want %d", len(buf.Data), tc.targetSize)
			}
		})
	}

}
