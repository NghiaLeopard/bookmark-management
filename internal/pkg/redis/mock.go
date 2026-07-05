package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func NewMockRClient(t *testing.T) *redis.Client {
	rdb := miniredis.RunT(t)

	return redis.NewClient(&redis.Options{
		Addr: rdb.Addr(),
	})
}
