package cache

import lru "github.com/hashicorp/golang-lru/v2"

type TwoQueue struct {
	v *lru.TwoQueueCache[string, string]
}

func NewTwoQueue(size int) Cache {
	v, err := lru.New2Q[string, string](size)
	if err != nil {
		panic(err)
	}

	return &TwoQueue{
		v: v,
	}
}

func (t *TwoQueue) Name() string {
	return "two-queue"
}

func (t *TwoQueue) Get(key string) bool {
	_, ok := t.v.Get(key)
	return ok
}

func (t *TwoQueue) Set(key string) {
	t.v.Add(key, key)
}

func (t *TwoQueue) Close() {

}
