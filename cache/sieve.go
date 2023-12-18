package cache

import (
	fifo "github.com/scalalang2/golang-fifo"
	sieve "github.com/scalalang2/golang-fifo/sieve"
)

type Sieve struct {
	v fifo.Cache[string, any]
}

func NewSieve(size int) Cache {
	return &Sieve{sieve.NewSieve[string, any](size)}
}

func (s *Sieve) Name() string {
	return "sieve"
}

func (s *Sieve) Get(key string) bool {
	_, ok := s.v.Get(key)
	return ok
}

func (s *Sieve) Set(key string) {
	s.v.Set(key, key)
}

func (s *Sieve) Close() {

}
