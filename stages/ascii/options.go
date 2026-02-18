package ascii

type opts struct {
	useEdge       bool
	edgeThreshold int
	invert        bool
	colorMode     ColorMode
}

type optFunc func(*opts)

func defaultOpts() opts {
	return opts{
		invert:        false,
		edgeThreshold: 20,
		useEdge:       true,
		colorMode:     FullColor,
	}
}

func WithEdge(useEdge bool) optFunc {
	return func(o *opts) {
		o.useEdge = useEdge
	}
}

func WithEdgeThreshold(threshold int) optFunc {
	return func(o *opts) {
		o.edgeThreshold = threshold
	}
}

func WithInvert(invert bool) optFunc {
	return func(o *opts) {
		o.invert = invert
	}
}

func WithColorMode(mode ColorMode) optFunc {
	return func(o *opts) {
		o.colorMode = mode
	}
}
