package cache

import (
	"github.com/Code-Hex/go-generics-cache/policy/clock"
	"sync"
)

type Clock struct {
	lock sync.Mutex
	v    *clock.Cache[string, string]
}

func NewClock(size int) Cache {
	return &Clock{
		v: clock.NewCache[string, string](clock.WithCapacity(size)),
	}
}

func (c *Clock) Name() string {
	return "clock"
}

func (c *Clock) Get(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.v.Get(key)
	return ok
}

func (c *Clock) Set(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.v.Set(key, key)
}

func (c *Clock) Close() {

}
