package main

import "github.com/NghiaLeopard/bookmark-management/internal/api"

func main() {
	app := api.NewEngine()

	app.Start()
}
