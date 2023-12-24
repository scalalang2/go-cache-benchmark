package cache

import (
	"github.com/maypok86/otter"
)

type Otter struct {
	v *otter.Cache[string, any]
}

func NewOtter(size int) Cache {
	cache, err := otter.MustBuilder[string, any](size).Build()
	if err != nil {
		panic(err)
	}

	return &Otter{v: cache}
}

func (c *Otter) Name() string {
	return "otter"
}

func (c *Otter) Get(key string) bool {
	_, ok := c.v.Get(key)
	return ok
}

func (c *Otter) Set(key string) {
	c.v.Set(key, key)
}

func (c *Otter) Close() {

}
