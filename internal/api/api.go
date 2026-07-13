package api

import (
	"log"
	"net/http"

	"github.com/NghiaLeopard/bookmark-management/docs"
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

	basePath := e.config.BasePath

	docs.SwaggerInfo.BasePath = basePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	healthCheckService := service.NewHealthCheck(e.config, e.rdb)
	healthCheckHandler := handler.NewHealthCheck(healthCheckService)

	genPassService := service.NewGenPassService()
	genPassHandler := handler.NewGenPassHandler(genPassService)

	urlStorage := repository.NewUrlStorage(e.rdb)
	urlService := service.NewShortenUrl(urlStorage, genPassService)
	urlHandler := handler.NewShortenUrlHandler(urlService)

	apiGroup := e.app.Group("v1")
	{
		apiGroup.GET("/health-check", healthCheckHandler.CheckHealth)
		apiGroup.POST("/genpass", genPassHandler.GeneratePassword)
		apiGroup.POST("/links/shorten", urlHandler.CreateShortenUrl)
		apiGroup.GET("/links/redirect/:code", urlHandler.Redirect)

	}

}
