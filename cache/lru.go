package cache

import lru "github.com/hashicorp/golang-lru/v2"

type LRU struct {
	v *lru.Cache[string, any]
}

func NewLRU(size int) Cache {
	v, _ := lru.New[string, any](size)
	return &LRU{v}
}

func (c *LRU) Name() string {
	return "lru-hashicorp"
}

func (c *LRU) Get(key string) bool {
	_, ok := c.v.Get(key)
	return ok
}

func (c *LRU) Set(key string) {
	c.v.Add(key, key)
}

func (c *LRU) Close() {

}
