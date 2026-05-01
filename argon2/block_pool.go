package argon2

import (
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
)

var (
	blockPoolMutex sync.Mutex
	blockPools     *lru.Cache[uint32, *sync.Pool]
)

func init() {
	poolsCache, err := lru.New[uint32, *sync.Pool](8)
	if err != nil {
		panic(fmt.Errorf("argon2: failed to create block pools cache: %w", err))
	}

	blockPools = poolsCache
}

func getOrCreateBlockPool(size uint32) *sync.Pool {
	if pool, ok := blockPools.Get(size); ok {
		return pool
	}

	return upsertBlockPool(size)
}

func upsertBlockPool(size uint32) *sync.Pool {
	blockPoolMutex.Lock()
	defer blockPoolMutex.Unlock()
	if pool, ok := blockPools.Get(size); ok {
		return pool
	}

	pool := &sync.Pool{
		New: func() any {
			b := make([]block, size)
			return &b
		},
	}

	blockPools.Add(size, pool)
	return pool
}
