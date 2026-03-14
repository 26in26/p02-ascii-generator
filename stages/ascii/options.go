package ascii

type stageOpts struct {
	invert bool
}

type stageOptFunc func(*stageOpts)

func defaultOpts() stageOpts {
	return stageOpts{
		invert: false,
	}
}

func WithInvert(invert bool) stageOptFunc {
	return func(o *stageOpts) {
		o.invert = invert
	}
}

type edgeFilterOpts struct {
	threshold int
}

type edgeFilterOptFunc func(*edgeFilterOpts)

func defaultEdgeOpts() edgeFilterOpts {
	return edgeFilterOpts{
		threshold: 20,
	}
}

func WithEdgeThreshold(threshold int) edgeFilterOptFunc {
	return func(o *edgeFilterOpts) {
		o.threshold = threshold
	}
}

type colorFilterOpts struct {
	colorMode ColorMode
}

type colorFilterOptFunc func(*colorFilterOpts)

func defaultColorOpts() colorFilterOpts {
	return colorFilterOpts{
		colorMode: FullColor,
	}
}

func WithColorMode(mode ColorMode) colorFilterOptFunc {
	return func(o *colorFilterOpts) {
		o.colorMode = mode
	}
}
