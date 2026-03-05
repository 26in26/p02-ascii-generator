package pipeline

import (
	"sync"

	"github.com/26in26/p02-ascii-generator/image"
	"github.com/26in26/p02-ascii-generator/utils"
)

type Frame struct {
	Id uint64
	wg sync.WaitGroup

	Raw Waitable[*image.RGBBuffer]

	ResizedImage Waitable[*image.RGBBuffer]
	GrayImage    Waitable[*image.GrayBuffer]
	GradientMap  Waitable[utils.Gradient]
	AsciiArt     Waitable[*image.AsciiBuffer]
}

func (f *Frame) Reset(id uint64, stages int) {
	f.Id = id
	f.wg.Add(stages)

	f.ResizedImage.Reset()
	f.GrayImage.Reset()
	f.GradientMap.Reset()
	f.AsciiArt.Reset()
}

type FramePool struct {
	pool *sync.Pool
}

func (p *FramePool) Get() *Frame {
	return p.pool.Get().(*Frame)
}

// Safely return the frame to the pool (waits till all the stages are done)
func (p *FramePool) Put(frame *Frame) {
	frame.wg.Wait()
	p.pool.Put(frame)
}

func NewFramePool() *FramePool {
	return &FramePool{
		pool: &sync.Pool{
			New: func() any {
				return &Frame{
					Raw: NewWaitable[*image.RGBBuffer](),

					ResizedImage: NewWaitable[*image.RGBBuffer](),
					GrayImage:    NewWaitable[*image.GrayBuffer](),
					GradientMap:  NewWaitable[utils.Gradient](),
					AsciiArt:     NewWaitable[*image.AsciiBuffer](),
				}
			},
		},
	}
}
