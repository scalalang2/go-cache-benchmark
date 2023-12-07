package cache

type Cache interface {
	Name() string
	Get(key string) bool
	Set(key string)
	Close()
}
