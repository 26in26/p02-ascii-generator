package flow

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/26in26/p02-ascii-generator/pipeline"
)

type packet[T any] struct {
	ID      uint64
	Data    T
	Err     error
	Release func()
}

type Outlet[T any] struct {
	ch chan *packet[T]
}

func NewOutlet[T any](size int) Outlet[T] {
	return Outlet[T]{
		ch: make(chan *packet[T], size),
	}
}

func (o *Outlet[T]) Feed(ctx context.Context, data <-chan T) {
	go func() {
		var id uint64
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-data:
				if !ok {
					return
				}
				o.ch <- &packet[T]{
					ID:      id,
					Data:    d,
					Release: func() {},
				}
				id++
			}
		}
	}()
}

func (o *Outlet[T]) Sink(f func(T, error)) {
	go func() {
		for p := range o.ch {
			f(p.Data, p.Err)
			p.Release()
		}
	}()
}

func (o *Outlet[T]) Branch(ctx context.Context, branches int) []Outlet[T] {
	outlets := make([]Outlet[T], branches)

	for i := range outlets {
		outlets[i] = NewOutlet[T](10)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case p := <-o.ch:
				var refs atomic.Int32
				refs.Store(int32(branches))

				// The "Virtual" release only triggers the real release when count is 0
				sharedRelease := func() {
					if refs.Add(-1) == 0 {
						p.Release()
					}
				}

				for _, outlet := range outlets {
					outlet.ch <- &packet[T]{
						ID:      p.ID,
						Data:    p.Data,
						Release: sharedRelease,
					}
				}
			}

		}

	}()

	return outlets
}

func (o *Outlet[T]) Close() {
	close(o.ch)
}

func Map[I, O any](ctx context.Context, in *Outlet[I], stage pipeline.Stage[I, O]) Outlet[O] {
	out := NewOutlet[O](10)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case p := <-in.ch:
				// if there is an error pass it on
				if p.Err != nil {
					out.ch <- &packet[O]{ID: p.ID, Err: p.Err, Release: func() {}}
					p.Release()
					continue
				}

				// TODO: Create sub context for self canclation
				data, err := stage.Kernal(ctx, p.Data)
				if err != nil {
					err = fmt.Errorf("stage %s: %w", stage.Name(), err)
				}

				out.ch <- &packet[O]{
					ID:      p.ID,
					Data:    data,
					Err:     err,
					Release: func() { stage.Release(data) },
				}
				p.Release()

			}
		}

	}()

	return out
}

func Mask[I, O any](ctx context.Context, in *Outlet[I], filter pipeline.Filter[I, O]) Outlet[O] {
	out := NewOutlet[O](10)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case p := <-in.ch:
				// if there is an error pass it on
				if p.Err != nil {
					out.ch <- &packet[O]{ID: p.ID, Err: p.Err, Release: func() {}}
					p.Release()
					continue
				}

				data, err := filter.Apply(ctx, p.Data)
				if err != nil {
					err = fmt.Errorf("filter %s: %w", filter.Name(), err)
				}

				out.ch <- &packet[O]{
					ID:      p.ID,
					Data:    data,
					Err:     err,
					Release: p.Release,
				}

			}
		}
	}()

	return out
}
