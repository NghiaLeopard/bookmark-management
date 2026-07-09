package redis

import (
	"github.com/redis/go-redis/v9"
)

func NewClient(prefix string) (*redis.Client, error) {
	cfg, err := newConfig(prefix)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.Db,
	})

	return rdb, nil
}
