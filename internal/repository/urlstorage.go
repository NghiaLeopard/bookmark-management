package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name UrlStorage --filename urlstorage.go
type UrlStorage interface {
	StoreUrl(ctx context.Context, code, url string, expire time.Duration) error
	GetUrlByCode(ctx context.Context, code string) (string, error)
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
	expire = expire * time.Second
	_, err := s.rdb.Set(ctx, code, url, expire).Result()
	return err
}

var ErrCodeNotFound = errors.New("code not found")

func (s *urlStorage) GetUrlByCode(ctx context.Context, code string) (string, error) {
	code, err := s.rdb.Get(ctx, code).Result()

	if errors.Is(err, redis.Nil) {
		return "", ErrCodeNotFound
	}

	if err != nil {
		return "", err
	}

	return code, nil
}
