package main

import (
	_ "github.com/NghiaLeopard/bookmark-management/docs"

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
