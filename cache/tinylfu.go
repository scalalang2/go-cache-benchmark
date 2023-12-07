package cache

import (
	"sync"

	"github.com/dgryski/go-tinylfu"
)

type TinyLFU[V any] struct {
	mu sync.Mutex
	v  *tinylfu.T[string]
}

func NewTinyLFU(size int) Cache {
	return &TinyLFU[any]{
		v: tinylfu.New[string](size, size*10),
	}
}

func (c *TinyLFU[V]) Name() string {
	return "tinylfu"
}

func (c *TinyLFU[V]) Set(key string) {
	c.mu.Lock()
	c.v.Add(key, key)
	c.mu.Unlock()
}

func (c *TinyLFU[V]) Get(key string) bool {
	c.mu.Lock()
	_, ok := c.v.Get(key)
	c.mu.Unlock()
	return ok
}

func (c *TinyLFU[V]) Close() {}
