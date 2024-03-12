package cache

import (
	lru "github.com/elastic/go-freelru"
	"github.com/zeebo/xxh3"
)

type FreeLRUSynced struct {
	v *lru.SyncedLRU[string, string]
}

func hash(s string) uint32 {
	return uint32(xxh3.HashString(s))
}

func NewFreeLRUSynced(size int) Cache {
	// The extra factor makes SyncedLRU use the same amount of memory as S3-FIFO does.
	v, _ := lru.NewSynced[string, string](uint32(size), hash)
	return &FreeLRUSynced{v}
}

func (c *FreeLRUSynced) Name() string {
	return "freelru-synced"
}

func (c *FreeLRUSynced) Get(key string) bool {
	_, ok := c.v.Get(key)
	return ok
}

func (c *FreeLRUSynced) Set(key string) {
	c.v.Add(key, key)
}

func (c *FreeLRUSynced) Close() {

}

type FreeLRUSharded struct {
	v *lru.ShardedLRU[string, string]
}

func NewFreeLRUSharded(size int) Cache {
	// The extra factor makes ShardedLRU use the same amount of memory as S3-FIFO does.
	v, _ := lru.NewShardedWithSize[string, string](128, uint32(size), uint32(size), hash)
	return &FreeLRUSharded{v}
}

func (c *FreeLRUSharded) Name() string {
	return "freelru-sharded"
}

func (c *FreeLRUSharded) Get(key string) bool {
	_, ok := c.v.Get(key)
	return ok
}

func (c *FreeLRUSharded) Set(key string) {
	c.v.Add(key, key)
}

func (c *FreeLRUSharded) Close() {

}
