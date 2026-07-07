package handler

import (
	"net/http"
	"time"

	"github.com/NghiaLeopard/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type ShortenUrlHandler interface {
	CreateShortenUrl(ctx *gin.Context)
	GetUrlByCode(ctx *gin.Context)
}

type shortenUrlHandler struct {
	shortenUrlService service.ShortenUrlService
}

func NewShortenUrlHandler(shortenUrlService service.ShortenUrlService) ShortenUrlHandler {
	return &shortenUrlHandler{
		shortenUrlService: shortenUrlService,
	}
}

type ShortenUrlInputBody struct {
	Url    string        `json:"url" binding:"required"`
	Expire time.Duration `json:"expire" binding:"required"`
}

func (h *shortenUrlHandler) CreateShortenUrl(ctx *gin.Context) {
	var input ShortenUrlInputBody
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Invalid input")
		return
	}

	code, err := h.shortenUrlService.CreateShortenUrl(ctx, input.Url, input.Expire)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "Internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": "Shorten URL generated successfully!",
	})
}

func (h *shortenUrlHandler) GetUrlByCode(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
	})
}
