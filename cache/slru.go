package cache

import (
	"github.com/golang/groupcache/lru"
	"sync"
)

type SLRU struct {
	lock  sync.Mutex
	once  *lru.Cache
	twice *lru.Cache
}

func NewSLRU(size int) Cache {
	return &SLRU{
		once:  lru.New(int(float64(size) * 0.2)),
		twice: lru.New(int(float64(size) * 0.8)),
	}
}

func (s *SLRU) Name() string {
	return "slru"
}

func (s *SLRU) Get(key string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	val, ok := s.once.Get(key)
	if ok {
		s.once.Remove(key)
		s.twice.Add(key, val)
		return true
	}

	_, ok = s.twice.Get(key)
	return ok
}

func (s *SLRU) Set(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.once.Add(key, key)
}

func (s *SLRU) Close() {
	s.once.Clear()
	s.twice.Clear()
}
