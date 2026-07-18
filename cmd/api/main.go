package main

import (
	_ "github.com/NghiaLeopard/bookmark-management/docs"

	"github.com/NghiaLeopard/bookmark-management/internal/api"
	"github.com/NghiaLeopard/bookmark-management/internal/pkg/redis"
)

// @title Bookmark Management API
// @version 1.0
// @description Bookmark management service API
func main() {
	rdb, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}
	defer rdb.Close()

	app := api.NewEngine(rdb)
	app.Start()
}
