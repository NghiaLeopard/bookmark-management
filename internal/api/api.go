package api

import (
	"net/http"

	"github.com/NghiaLeopard/bookmark-management/internal/handler"
	"github.com/NghiaLeopard/bookmark-management/internal/repository"
	"github.com/NghiaLeopard/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Engine interface {
	Start()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app *gin.Engine
	rdb *redis.Client
}

func NewEngine(rdb *redis.Client) Engine {
	app := &engine{app: gin.Default(), rdb: rdb}

	app.InitRoutes()

	return app
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

func (e *engine) Start() {
	e.app.Run(":8080")
}

func (e *engine) InitRoutes() {
	genPassService := service.NewGenPassService()
	genPassHandler := handler.NewGenPassHandler(genPassService)

	urlStorage := repository.NewUrlStorage(e.rdb)
	urlService := service.NewShortenUrlService(urlStorage, genPassService)
	urlHandler := handler.NewShortenUrlHandler(urlService)
	e.app.POST("/genpass", genPassHandler.GeneratePassword)
	e.app.POST("/shortenurl", urlHandler.CreateShortenUrl)
}
