package cache

import (
	"context"
	"log"
	"sync"
	"time"

	goCache "github.com/patrickmn/go-cache"

	"service/internal/storage/postgres"
)

type Cache struct {
	*goCache.Cache
}

func New(ctx context.Context, storage *postgres.Storage, wg *sync.WaitGroup) (*Cache, error) {
	const op = "internal.cache.New"

	var cache Cache
	cache.Cache = goCache.New(1*time.Hour, 24*time.Hour)

	all, err := storage.GetAll()
	if err != nil {
		return &cache, err
	}

	wg.Add(1)
	go func(all map[string][]byte, cache *Cache, wg *sync.WaitGroup) {
		defer log.Println("Add all orders to cache successful!")
		defer wg.Done()
		for k, v := range all {
			select {
			case <-ctx.Done():
				log.Println("Cache canceled by context!")
				return
			default:
				cache.SetDefault(k, v)
			}
		}

	}(all, &cache, wg)

	return &cache, nil
}
