package cache

import (
	fifo "github.com/scalalang2/golang-fifo"
	"github.com/scalalang2/golang-fifo/s3fifo"
)

type S3FIFO struct {
	v fifo.Cache[string, any]
}

func NewS3FIFO(size int) Cache {
	v := s3fifo.New[string, any](size)
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
