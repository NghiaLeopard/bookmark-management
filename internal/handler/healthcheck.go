package handler

import (
	"net/http"

	"github.com/NghiaLeopard/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type HealthCheck interface {
	CheckHealth(c *gin.Context)
}

type healthCheckHandler struct {
	service service.HealthCheck
}

func NewHealthCheck(service service.HealthCheck) HealthCheck {
	return &healthCheckHandler{
		service: service,
	}
}

// CheckHealth godoc
// @Summary check health of the service
// @Schemes
// @Description check health of the service
// @Tags health-check
// @Accept json
// @Produce json
// @Success 200 {object} model.HealthCheck
// @Router /health-check [get]
func (h *healthCheckHandler) CheckHealth(c *gin.Context) {
	healthCheck, err := h.service.CheckHealth()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, healthCheck)
}
