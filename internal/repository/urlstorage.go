package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type UrlStorage interface {
	StoreUrl(ctx context.Context, code, url string, expire time.Duration) error
}

type urlStorage struct {
	rdb *redis.Client
}

func NewUrlStorage(rdb *redis.Client) UrlStorage {
	return &urlStorage{
		rdb: rdb,
	}
}

func (s *urlStorage) StoreUrl(ctx context.Context, code, url string, expire time.Duration) error {
	_, err := s.rdb.Set(ctx, code, url, expire).Result()
	return err
}

func (s *urlStorage) GetUrl(ctx context.Context, code string) (string, error) {
	return s.rdb.Get(ctx, code).Result()
}
