package charcolor

import (
	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/internal/impl/filters"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/pipeline/flow"
	"github.com/26in26/p02-ascii-generator/utils/builder"
)

var FullColor = filters.FullColor
var RetroColor = filters.RetroColor
var Monochrome = filters.Monochrome

type colorOpts struct {
	colorMode filters.ColorMode
}

type colorOptFunc = builder.Option[colorOpts]

func defaultColorOpts() *colorOpts {
	return &colorOpts{
		colorMode: FullColor,
	}
}

func WithColorMode(mode filters.ColorMode) colorOptFunc {
	return func(o *colorOpts) {
		o.colorMode = mode
	}
}

func NewColorFilter(opts ...colorOptFunc) pipeline.Filter[flow.Pair[*image.AsciiBuffer, *image.RGBBuffer], *image.AsciiBuffer] {
	cfg := builder.Build(defaultColorOpts, opts...)

	return filters.NewColorFilter(cfg.colorMode)
}
