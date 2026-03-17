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
