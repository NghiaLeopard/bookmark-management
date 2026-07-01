package handler

import (
	"net/http"

	"github.com/NghiaLeopard/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type GenPass interface {
	GeneratePassword(c *gin.Context)
}

type genPassHandler struct {
	genPassService service.GenPass
}

func NewGenPassHandler(genPassService service.GenPass) GenPass {
	return &genPassHandler{genPassService: genPassService}
}

var lengthPassword int = 12

func (h *genPassHandler) GeneratePassword(c *gin.Context) {

	password, err := h.genPassService.GeneratePassword(lengthPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"password": password})

}
