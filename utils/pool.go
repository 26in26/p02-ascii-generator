package utils

import "sync"

// Pool is a generic wrapper around sync.Map and sync.Pool.
// It allows managing multiple sync.Pools keyed by a comparable type K,
// storing values of type V.
type Pool[K comparable, V any] struct {
	pools   sync.Map // map[K]*sync.Pool
	factory func(K) (V, error)
}

// NewPool creates a new Pool.
// factory is a function that creates a new instance of V given a key K.
func NewPool[K comparable, V any](factory func(K) (V, error)) *Pool[K, V] {
	return &Pool[K, V]{
		factory: factory,
	}
}

// Get retrieves an item from the pool specific to the given key.
// If the specific pool is empty or doesn't exist, a new item is created using the factory.
func (p *Pool[K, V]) Get(key K) (V, error) {
	pool := p.getPool(key)
	item := pool.Get()

	if item == nil {
		return p.factory(key)
	}

	return item.(V), nil
}

// Put returns an item to the pool specific to the given key.
func (p *Pool[K, V]) Put(key K, item V) {
	p.getPool(key).Put(item)
}

func (p *Pool[K, V]) getPool(key K) *sync.Pool {
	if v, ok := p.pools.Load(key); ok {
		return v.(*sync.Pool)
	}

	newPool := &sync.Pool{
		New: func() any {
			return nil
		},
	}

	actual, _ := p.pools.LoadOrStore(key, newPool)
	return actual.(*sync.Pool)
}
