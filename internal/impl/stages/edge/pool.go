package edge

import (
	"github.com/26in26/p02-ascii-generator/utils"
	"github.com/26in26/p02-ascii-generator/utils/pools"
)

type gradientPool struct {
	*pools.KeyedPool[int, utils.Gradient]
}

func NewGrayImagePool() *gradientPool {
	return &gradientPool{
		KeyedPool: pools.NewPool(func(s int) (utils.Gradient, error) {
			return make(utils.Gradient, s), nil
		}),
	}
}
