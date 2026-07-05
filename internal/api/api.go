package api

import (
	"log"
	"net/http"

	"github.com/NghiaLeopard/bookmark-management/internal/config"
	"github.com/NghiaLeopard/bookmark-management/internal/handler"
	"github.com/NghiaLeopard/bookmark-management/internal/repository"
	"github.com/NghiaLeopard/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Engine interface {
	Start()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app    *gin.Engine
	config *config.Config
	rdb    *redis.Client
}

func NewEngine(rdb *redis.Client) Engine {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	app := &engine{app: gin.Default(), config: cfg, rdb: rdb}
	app.InitRoutes()

	return app
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

func (e *engine) Start() {
	e.app.Run(":" + e.config.Port)
}

func (e *engine) InitRoutes() {
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	healthCheckService := service.NewHealthCheck(e.config)
	healthCheckHandler := handler.NewHealthCheck(healthCheckService)
	e.app.GET("/health-check", healthCheckHandler.CheckHealth)

	genPassService := service.NewGenPassService()
	genPassHandler := handler.NewGenPassHandler(genPassService)
	e.app.POST("/genpass", genPassHandler.GeneratePassword)

	if e.rdb != nil {
		urlStorage := repository.NewUrlStorage(e.rdb)
		urlService := service.NewShortenUrlService(urlStorage, genPassService)
		urlHandler := handler.NewShortenUrlHandler(urlService)
		e.app.POST("/shortenurl", urlHandler.CreateShortenUrl)
	}
}
