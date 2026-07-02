package main

import (
	_ "github.com/NghiaLeopard/bookmark-management/docs"

	"github.com/NghiaLeopard/bookmark-management/internal/api"
)

func main() {
	app := api.NewEngine()

	app.Start()
}
