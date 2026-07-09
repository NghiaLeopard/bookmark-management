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
	Url    string        `json:"url" binding:"required,url" example:"https://example.com/long-page"`
	Expire time.Duration `json:"expire" binding:"required,min=0" swaggertype:"integer" example:"3600000000000"`
}

// CreateShortenUrl godoc
// @Summary Create shorten url
// @Schemes
// @Description Create a short code for a long URL and store it in Redis until expire (duration in nanoseconds, e.g. 3600000000000 = 1 hour)
// @Tags links
// @Accept json
// @Produce json
// @Param body body ShortenUrlInputBody true "URL to shorten and TTL"
// @Success 200 {object} model.ShortenUrlResponse
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /links/shorten [post]
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
