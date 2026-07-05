package main

import (
	"github.com/NghiaLeopard/bookmark-management/internal/api"
	"github.com/NghiaLeopard/bookmark-management/internal/pkg/redis"
)

func main() {

	rdb, err := redis.NewClient("")

	if err != nil {
		panic(err)
	}

	defer rdb.Close()
	app := api.NewEngine(rdb)

	app.Start()
}
