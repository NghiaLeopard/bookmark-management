package api

import (
	"log"
	"net/http"

	"github.com/NghiaLeopard/bookmark-management/internal/config"
	"github.com/NghiaLeopard/bookmark-management/internal/handler"
	"github.com/NghiaLeopard/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
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
}

func NewEngine() Engine {
	config, err := config.NewConfig()

	if err != nil {
		log.Fatal(err.Error())
	}
	app := &engine{app: gin.Default(), config: config}

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
}
