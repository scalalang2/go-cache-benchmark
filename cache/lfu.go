package cache

import (
	"github.com/Code-Hex/go-generics-cache/policy/lfu"
)

type LFU struct {
	c *lfu.Cache[string, string]
}

func NewLFU(size int) Cache {
	return &LFU{
		c: lfu.NewCache[string, string](lfu.WithCapacity(size)),
	}
}

func (c *LFU) Name() string {
	return "lfu"
}

func (c *LFU) Get(key string) bool {
	_, ok := c.c.Get(key)
	return ok
}

func (c *LFU) Set(key string) {
	c.c.Set(key, key)
}

func (c *LFU) Close() {

}
