package api

import (
	"net/http"

	"github.com/NghiaLeopard/bookmark-management/internal/handler"
	"github.com/NghiaLeopard/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type Engine interface {
	Start()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app *gin.Engine
}

func NewEngine() Engine {
	app := &engine{app: gin.Default()}

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
	e.app.POST("/genpass", genPassHandler.GeneratePassword)
}
