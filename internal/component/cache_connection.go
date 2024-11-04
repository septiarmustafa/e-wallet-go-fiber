package component

import (
	"context"
	"e-wallet/domain"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
)

func GetCacheConnection() domain.CacheRepository {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		log.Fatalf("error when connect cache %s", err.Error())
	}
	return cache
}
