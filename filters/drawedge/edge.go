package drawedge

import (
	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/internal/impl/filters"
	"github.com/26in26/p02-ascii-generator/pipeline"
	"github.com/26in26/p02-ascii-generator/pipeline/flow"
	"github.com/26in26/p02-ascii-generator/utils"
	"github.com/26in26/p02-ascii-generator/utils/builder"
)

type edgeFilterOpts struct {
	threshold int
}

type edgeOptFunc = builder.Option[edgeFilterOpts]

func defaultEdgeOpts() *edgeFilterOpts {
	return &edgeFilterOpts{
		threshold: 20,
	}
}

func WithEdgeThreshold(threshold int) edgeOptFunc {
	return func(o *edgeFilterOpts) {
		o.threshold = threshold
	}
}

func NewEdgeFilter(opts ...edgeOptFunc) pipeline.Filter[flow.Pair[*image.AsciiBuffer, utils.Gradient], *image.AsciiBuffer] {
	cfg := builder.Build(defaultEdgeOpts, opts...)

	return filters.NewEdgeFilter(cfg.threshold)
}
