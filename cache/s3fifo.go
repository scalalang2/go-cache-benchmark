package cache

import (
	"github.com/scalalang2/golang-fifo/s3fifo"
	"github.com/scalalang2/golang-fifo/types"
)

type S3FIFO struct {
	v types.Cache[string, any]
}

func NewS3FIFO(size int) Cache {
	v := s3fifo.New[string, any](size, 0)
	return &S3FIFO{v}
}

func (c *S3FIFO) Name() string {
	return "s3-fifo"
}

func (c *S3FIFO) Get(key string) bool {
	_, ok := c.v.Get(key)
	return ok
}

func (c *S3FIFO) Set(key string) {
	c.v.Set(key, key)
}

func (c *S3FIFO) Close() {

}
