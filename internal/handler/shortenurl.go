package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/NghiaLeopard/bookmark-management/internal/repository"
	"github.com/NghiaLeopard/bookmark-management/internal/service"
	"github.com/NghiaLeopard/bookmark-management/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ShortenUrlHandler interface {
	CreateShortenUrl(ctx *gin.Context)
	Redirect(ctx *gin.Context)
}

type shortenUrlHandler struct {
	shortenUrlService service.ShortenUrl
}

func NewShortenUrlHandler(shortenUrlService service.ShortenUrl) ShortenUrlHandler {
	return &shortenUrlHandler{
		shortenUrlService: shortenUrlService,
	}
}

type ShortenUrlInputBody struct {
	Url    string        `json:"url" binding:"required,url" example:"https://example.com/long-page"`
	Expire time.Duration `json:"exp" binding:"required,min=0" swaggertype:"integer" example:"3600"`
}

// CreateShortenUrl godoc
// @Summary Create shorten url
// @Schemes
// @Description Create a short code for a long URL and store it in Redis until expire (duration in seconds, e.g. 3600 = 1 hour)
// @Tags links
// @Accept json
// @Produce json
// @Param body body ShortenUrlInputBody true "URL to shorten and TTL"
// @Success 200 {object} model.ShortenUrlResponse
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/links/shorten [post]
func (h *shortenUrlHandler) CreateShortenUrl(ctx *gin.Context) {
	var input ShortenUrlInputBody
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	code, err := h.shortenUrlService.CreateShortenUrl(ctx, input.Url, input.Expire)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.shortenurl.CreateShortenUrl").Msg("Failed to create shorten url")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": "Shorten URL generated successfully!",
	})
}

// Redirect godoc
// @Summary Redirect to long URL
// @Schemes
// @Description Redirect to long URL by code
// @Tags links
// @Param code path string true "Code"
// @Success 302 {string} string "Redirect to long URL"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Code not found"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/links/redirect/{code} [get]
func (h *shortenUrlHandler) Redirect(ctx *gin.Context) {
	code := ctx.Param("code")

	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputErrResponse)
		return
	}

	url, err := h.shortenUrlService.GetUrlByCode(ctx, code)

	if err != nil {
		if errors.Is(err, repository.ErrCodeNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, response.Message{
				Message: "Code not found",
				Detail:  nil,
			})
			return
		}

		log.Error().Err(err).Str("from", "handler.shortenurl.Redirect").Msg("Failed to get URL by code")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	ctx.Redirect(http.StatusFound, url)
}
