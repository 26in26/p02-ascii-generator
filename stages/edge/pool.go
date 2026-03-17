package edge

import (
	"github.com/26in26/p02-ascii-generator/utils"
)

type gradientPool struct {
	*utils.Pool[int, utils.Gradient]
}

func NewGrayImagePool() *gradientPool {
	return &gradientPool{
		Pool: utils.NewPool(func(s int) (utils.Gradient, error) {
			return make(utils.Gradient, s), nil
		}),
	}
}
