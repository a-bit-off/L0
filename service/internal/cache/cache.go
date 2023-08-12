package cache

import (
	"log"
	"sync"
	"time"

	goCache "github.com/patrickmn/go-cache"

	"service/internal/storage/postgres"
)

type Cache struct {
	*goCache.Cache
}

func New(storage *postgres.Storage, wg *sync.WaitGroup) (*Cache, error) {
	var cache Cache
	var err error

	cache.Cache = goCache.New(1*time.Hour, 24*time.Hour)

	wg.Add(1)
	go func(cache Cache, err error) {
		wg.Done()
		all, err := storage.GetAll()
		if err != nil {
			return
		}

		for k, v := range all {
			cache.SetDefault(k, v)
		}

		log.Println("Add all orders to cache successful!")
	}(cache, err)

	return &cache, nil
}
