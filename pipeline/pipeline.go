package pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/26in26/p02-ascii-generator/image"
)

type Pipeline struct {
	stages []Stage
	pool   *FramePool
	nextId uint64

	results chan *Frame
	errors  chan error
}

func NewPipeline(stages ...Stage) (*Pipeline, error) {
	p :=
		&Pipeline{
			stages:  stages,
			pool:    NewFramePool(),
			nextId:  0,
			results: make(chan *Frame, 10),
			errors:  make(chan error, 10),
		}

	err := p.validate()

	return p, err
}

func (p *Pipeline) Enqueue(ctx context.Context, img *image.RGBBuffer) {
	f := p.pool.Get()
	f.Reset(p.nextId, len(p.stages))
	p.nextId++

	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// Async cleanup
	go func() {
		defer cancel()

		select {
		case <-ctx.Done():
			p.errors <- ctx.Err()
		case <-f.AsciiArt.Ready():
			p.results <- f
		}

		p.pool.Put(f)
	}()

	// goroutines
	for _, stage := range p.stages {
		go func(s Stage) {
			defer f.wg.Done()
			if err := s.Process(ctx, f); err != nil {
				p.errors <- fmt.Errorf("stage %s: %w", s.Name(), err)
			}

		}(stage)
	}

	f.Raw.Set(img)

}

func (p *Pipeline) Stream(ctx context.Context, src <-chan *image.RGBBuffer) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case img, ok := <-src:
				if !ok {
					return
				}
				p.Enqueue(ctx, img)
			}
		}
	}()
}

func (p *Pipeline) validate() error {
	return nil
}

func (p *Pipeline) Results() <-chan *Frame {
	return p.results
}

func (p *Pipeline) Errors() <-chan error {
	return p.errors
}
