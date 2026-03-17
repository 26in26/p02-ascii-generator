package resize

import (
	"sync"

	"github.com/26in26/p02-ascii-generator/image"
)

type RGBImagePool struct {
	pool *sync.Pool
	w    int
	h    int
}

func NewRGBImagePool(w, h int) *RGBImagePool {
	return &RGBImagePool{
		pool: &sync.Pool{
			New: func() any {
				return nil
			},
		},
		w: w,
		h: h,
	}
}

func (p *RGBImagePool) Get() (*image.RGBBuffer, error) {
	var err error
	img := p.pool.Get()
	if img == nil {
		img, err = image.NewRGBBuffer(p.w, p.h)
	}

	return img.(*image.RGBBuffer), err

}

func (p *RGBImagePool) Put(img *image.RGBBuffer) {
	if img == nil {
		return
	}

	p.pool.Put(img)
}

type xOffsetsPool struct {
	pool *sync.Pool
	len  int
}

func newXOffsetsPool(len int) *xOffsetsPool {
	return &xOffsetsPool{
		pool: &sync.Pool{
			New: func() any {
				return nil
			},
		},
		len: len,
	}
}

func (p *xOffsetsPool) Get() []int {

	offsets := p.pool.Get()
	if offsets == nil {
		offsets = make([]int, p.len)
	}

	return offsets.([]int)
}

func (p *xOffsetsPool) Put(offsets []int) {
	if offsets == nil {
		return
	}

	p.pool.Put(offsets)
}
