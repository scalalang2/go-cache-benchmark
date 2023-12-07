package cache

import "github.com/golang/groupcache/lru"

type LRUGroupCache struct {
	v *lru.Cache
}

func NewLRUGroupCache(size int) Cache {
	return &LRUGroupCache{
		v: lru.New(size),
	}
}

func (c *LRUGroupCache) Name() string {
	return "lru-groupcache"
}

func (c *LRUGroupCache) Get(key string) bool {
	_, ok := c.v.Get(key)
	return ok
}

func (c *LRUGroupCache) Set(key string) {
	c.v.Add(key, key)
}

func (c *LRUGroupCache) Close() {
	c.v.Clear()
}
