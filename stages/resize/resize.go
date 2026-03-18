package resize

import (
	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/internal/impl/stages/resize"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/utils/builder"
)

const CHAR_ASPECT_RATIO = 2

type opts struct {
	width  int
	height int
}

type optFunc = builder.Option[opts]

func defaultOpts() *opts {
	return &opts{
		width:  100,
		height: 100,
	}
}

func With(w, h int) optFunc {
	return func(o *opts) {
		o.width = w
		o.height = h
	}
}

func WithWidth(width int) optFunc {
	return func(o *opts) {
		o.width = width
	}
}

func WithHeight(height int) optFunc {
	return func(o *opts) {
		o.height = height
	}
}

func WithAspectRatio(w, h int, preserveWidth bool) optFunc {
	return func(o *opts) {
		aspectRatio := float64(w) / float64(h) * CHAR_ASPECT_RATIO
		if preserveWidth {
			// adjust height according to aspect ratio
			o.height = int(float64(o.width) / aspectRatio)
			if o.height <= 0 {
				o.height = 1
			}
		} else {
			// adjust width according to aspect ratio
			o.width = int(float64(o.height) * aspectRatio)
			if o.width <= 0 {
				o.width = 1
			}
		}
	}
}

func NewResizeStage(opts ...optFunc) (pipeline.Stage[*image.RGBBuffer, *image.RGBBuffer], error) {
	cfg := builder.Build(defaultOpts, opts...)

	return resize.NewResizeStage(cfg.width, cfg.height)
}
