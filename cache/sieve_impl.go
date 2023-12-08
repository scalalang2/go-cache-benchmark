package cache

type SieveCache struct {
	v *Sieve[string, any]
}

func NewSieveCache(size int) Cache {
	v := NewSieve[string, any](size)
	return &SieveCache{v}
}

func (c *SieveCache) Name() string {
	return "sieve"
}

func (c *SieveCache) Get(key string) bool {
	_, ok := c.v.Get(key)
	return ok
}

func (c *SieveCache) Set(key string) {
	c.v.Set(key, key)
}

func (c *SieveCache) Close() {

}
