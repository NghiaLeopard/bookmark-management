package main

import (
	"context"
	"time"

	"github.com/NghiaLeopard/bookmark-management/internal/pkg/redis"
)

func main() {
	rClient, err := redis.NewClient("")

	if err != nil {
		panic(err)
	}

	rClient.Set(context.Background(), "test", "test", time.Hour)

	rClient2, err := redis.NewClient("CACHE")

	if err != nil {
		panic(err)
	}

	rClient2.Set(context.Background(), "test-1", "test", time.Hour)

}
