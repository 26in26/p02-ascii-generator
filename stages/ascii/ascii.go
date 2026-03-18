package ascii

import (
	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/internal/impl/stages/ascii"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/utils/builder"
)

const DEFAULT_CHARSET = " .-=+*x#$&X@"

type asciiOpts struct {
	invert         bool
	densityCharset []byte
}

type optFunc = builder.Option[asciiOpts]

func defaultOpts() *asciiOpts {
	return &asciiOpts{
		invert:         false,
		densityCharset: []byte(DEFAULT_CHARSET),
	}
}

func WithInvert(invert bool) optFunc {
	return func(o *asciiOpts) {
		o.invert = invert
	}
}

type DensityCharset byte

const (
	StandardCharset DensityCharset = iota
	DenseCharset
	DotsCharset
)

func WithDensityCharset(charset DensityCharset) optFunc {
	return func(o *asciiOpts) {
		switch charset {
		case StandardCharset:
			o.densityCharset = []byte(DEFAULT_CHARSET)
		case DenseCharset:
			o.densityCharset = []byte(" .-=+*x#$&X@")
		case DotsCharset:
			o.densityCharset = []byte("  ·•")
		}

	}
}

func WithCustomCharset(charset string) optFunc {
	return func(o *asciiOpts) {
		o.densityCharset = []byte(charset)
	}
}

func NewAsciiStage(opts ...optFunc) pipeline.Stage[*image.GrayBuffer, *image.AsciiBuffer] {
	cfg := builder.Build(defaultOpts, opts...)

	return ascii.NewAsciiStage(cfg.densityCharset, cfg.invert)
}
