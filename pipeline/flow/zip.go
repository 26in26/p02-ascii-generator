package flow

import (
	"context"
	"errors"
)

type Pair[T, V any] struct {
	A T
	B V
}

func Zip[A, B any](ctx context.Context, inA Outlet[A], inB Outlet[B]) *Outlet[Pair[A, B]] {
	out := Outlet[Pair[A, B]]{ch: make(chan *packet[Pair[A, B]], 10)}

	go func() {
		pendingA := make(map[uint64]*packet[A])
		pendingB := make(map[uint64]*packet[B])

		for {
			select {
			case <-ctx.Done():
				return
			case pA := <-inA.ch:
				handlePairPacketA(pA, pendingA, pendingB, out)
			case pB := <-inB.ch:
				handlePairPacketB(pB, pendingA, pendingB, out)
			}
		}
	}()
	return &out
}

func handlePairPacketA[A, B any](pA *packet[A], pendingA map[uint64]*packet[A], pendingB map[uint64]*packet[B], out Outlet[Pair[A, B]]) {
	if pB, ok := pendingB[pA.ID]; ok {
		delete(pendingB, pB.ID)
		emitZip(pA, pB, out)
		return
	}
	pendingA[pA.ID] = pA
}

func handlePairPacketB[A, B any](pB *packet[B], pendingA map[uint64]*packet[A], pendingB map[uint64]*packet[B], out Outlet[Pair[A, B]]) {
	if pA, ok := pendingA[pB.ID]; ok {
		delete(pendingA, pA.ID)
		emitZip(pA, pB, out)
		return
	}
	pendingB[pB.ID] = pB
}

func emitZip[A, B any](pA *packet[A], pB *packet[B], out Outlet[Pair[A, B]]) {

	out.ch <- &packet[Pair[A, B]]{
		ID: pA.ID,
		Data: Pair[A, B]{
			A: pA.Data,
			B: pB.Data,
		},
		Err:     errors.Join(pA.Err, pB.Err),
		Release: func() { pA.Release(); pB.Release() },
	}

}

type Triplet[T, V, W any] struct {
	A T
	B V
	C W
}

func Zip3[A, B, C any](ctx context.Context, inA Outlet[A], inB Outlet[B], inC Outlet[C]) *Outlet[Triplet[A, B, C]] {
	out := NewOutlet[Triplet[A, B, C]](10)

	go func() {
		pendingA := make(map[uint64]*packet[A])
		pendingB := make(map[uint64]*packet[B])
		pendingC := make(map[uint64]*packet[C])

		for {
			select {
			case <-ctx.Done():
				return
			case pA := <-inA.ch:
				handleTripletPacketA(pA, pendingA, pendingB, pendingC, out)
			case pB := <-inB.ch:
				handleTripletPacketB(pB, pendingA, pendingB, pendingC, out)
			case pC := <-inC.ch:
				handleTripletPacketC(pC, pendingA, pendingB, pendingC, out)
			}
		}
	}()

	return &out
}

func handleTripletPacketA[A, B, C any](pA *packet[A], pendingA map[uint64]*packet[A], pendingB map[uint64]*packet[B], pendingC map[uint64]*packet[C], out Outlet[Triplet[A, B, C]]) {
	if pB, okB := pendingB[pA.ID]; okB {
		if pC, okC := pendingC[pA.ID]; okC {
			delete(pendingB, pB.ID)
			delete(pendingC, pC.ID)
			emitZip3(pA, pB, pC, out)
			return
		}
	}
	pendingA[pA.ID] = pA
}

func handleTripletPacketB[A, B, C any](pB *packet[B], pendingA map[uint64]*packet[A], pendingB map[uint64]*packet[B], pendingC map[uint64]*packet[C], out Outlet[Triplet[A, B, C]]) {
	if pA, okA := pendingA[pB.ID]; okA {
		if pC, okC := pendingC[pB.ID]; okC {
			delete(pendingA, pA.ID)
			delete(pendingC, pC.ID)
			emitZip3(pA, pB, pC, out)
			return
		}
	}
	pendingB[pB.ID] = pB
}

func handleTripletPacketC[A, B, C any](pC *packet[C], pendingA map[uint64]*packet[A], pendingB map[uint64]*packet[B], pendingC map[uint64]*packet[C], out Outlet[Triplet[A, B, C]]) {
	if pA, okA := pendingA[pC.ID]; okA {
		if pB, okB := pendingB[pC.ID]; okB {
			delete(pendingA, pA.ID)
			delete(pendingB, pB.ID)
			emitZip3(pA, pB, pC, out)
			return
		}
	}
	pendingC[pC.ID] = pC
}

func emitZip3[A, B, C any](pA *packet[A], pB *packet[B], pC *packet[C], out Outlet[Triplet[A, B, C]]) {

	out.ch <- &packet[Triplet[A, B, C]]{
		ID: pA.ID,
		Data: Triplet[A, B, C]{
			A: pA.Data,
			B: pB.Data,
			C: pC.Data,
		},
		Err:     errors.Join(pA.Err, pB.Err, pC.Err),
		Release: func() { pA.Release(); pB.Release(); pC.Release() },
	}

}
