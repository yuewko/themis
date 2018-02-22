package pip

import (
	"github.com/hashicorp/golang-lru"
)

type QueryCache struct {
	cache *lru.TwoQueueCache
	size  int
}

func NewQueryCache(size int) (*QueryCache, error) {
	cache, err := lru.New2Q(size)
	if err != nil {
		return nil, err
	}

	return &QueryCache{cache: cache, size: size}, nil
}

func (qc *QueryCache) Get(key interface{}) (interface{}, bool) {
	return qc.cache.Get(key)
}

func (qc *QueryCache) Add(key interface{}, value interface{}) {
	qc.cache.Add(key, value)
}

func (qc *QueryCache) Purge() {
	qc.cache.Purge()
}

func (qc *QueryCache) Len() int {
	return qc.cache.Len()
}
