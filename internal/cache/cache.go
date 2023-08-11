package cache

import (
	goCache "github.com/patrickmn/go-cache"
	"sync"
	"time"

	"L0/internal/storage/postgres"
)

type Cache struct {
	*goCache.Cache
}

func New(storage *postgres.Storage, wg *sync.WaitGroup) (*Cache, error) {
	var cache Cache
	cache.Cache = goCache.New(1*time.Hour, 24*time.Hour)

	all, err := storage.GetAll()
	if err != nil {
		return &cache, err
	}
	wg.Add(1)
	go func(all map[string][]byte, cache *Cache, wg *sync.WaitGroup) {
		defer wg.Done()
		for k, v := range all {
			cache.SetDefault(k, v)
		}

	}(all, &cache, wg)

	return &cache, nil
}
