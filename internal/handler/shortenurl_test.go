package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/NghiaLeopard/bookmark-management/internal/repository"
	mocksService "github.com/NghiaLeopard/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		Body               ShortenUrlInputBody
		SetUpRequest       func(ctx *gin.Context, body ShortenUrlInputBody)
		SetupMocksServices func(t *testing.T) *mocksService.ShortenUrlService
		ExpectedStatusCode int
		ExpectedResponse   string
	}{
		{
			name: "success",
			Body: ShortenUrlInputBody{
				Url:    "http://localhost:8000",
				Expire: time.Hour,
			},
			SetUpRequest: func(ctx *gin.Context, body ShortenUrlInputBody) {
				raw, err := json.Marshal(body)
				if err != nil {
					t.Fatalf("marshal body: %v", err)
				}

				ctx.Request = httptest.NewRequest(http.MethodPost, "/shorten-url", strings.NewReader(string(raw)))
			},
			SetupMocksServices: func(t *testing.T) *mocksService.ShortenUrlService {
				mockShortenUrlService := mocksService.NewShortenUrlService(t)
				mockShortenUrlService.On("CreateShortenUrl", mock.Anything, "http://localhost:8000", time.Hour).Return("1234567", nil)
				return mockShortenUrlService
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"code":"1234567","message":"Shorten URL generated successfully!"}`,
		},
		{
			name: "Invalid input",
			Body: ShortenUrlInputBody{
				Expire: time.Hour,
			},
			SetUpRequest: func(ctx *gin.Context, body ShortenUrlInputBody) {
				raw, err := json.Marshal(body)
				if err != nil {
					t.Fatalf("marshal body: %v", err)
				}

				ctx.Request = httptest.NewRequest(http.MethodPost, "/shorten-url", strings.NewReader(string(raw)))
			},
			SetupMocksServices: func(t *testing.T) *mocksService.ShortenUrlService {
				mockShortenUrlService := mocksService.NewShortenUrlService(t)
				return mockShortenUrlService
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `"Invalid input"`,
		},
		{
			name: "Invalid input",
			Body: ShortenUrlInputBody{
				Url: "http://localhost:8000",
			},
			SetUpRequest: func(ctx *gin.Context, body ShortenUrlInputBody) {
				raw, err := json.Marshal(body)
				if err != nil {
					t.Fatalf("marshal body: %v", err)
				}

				ctx.Request = httptest.NewRequest(http.MethodPost, "/shorten-url", strings.NewReader(string(raw)))
			},
			SetupMocksServices: func(t *testing.T) *mocksService.ShortenUrlService {
				mockShortenUrlService := mocksService.NewShortenUrlService(t)
				return mockShortenUrlService
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `"Invalid input"`,
		},
		{
			name: "Internal server error",
			Body: ShortenUrlInputBody{
				Url:    "http://localhost:8000",
				Expire: time.Hour,
			},
			SetUpRequest: func(ctx *gin.Context, body ShortenUrlInputBody) {
				raw, err := json.Marshal(body)
				if err != nil {
					t.Fatalf("marshal body: %v", err)
				}

				ctx.Request = httptest.NewRequest(http.MethodPost, "/shorten-url", strings.NewReader(string(raw)))
			},
			SetupMocksServices: func(t *testing.T) *mocksService.ShortenUrlService {
				mockShortenUrlService := mocksService.NewShortenUrlService(t)
				mockShortenUrlService.On("CreateShortenUrl", mock.Anything, "http://localhost:8000", time.Hour).Return("", errors.New("Internal server error"))
				return mockShortenUrlService
			},
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedResponse:   `"Internal server error"`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			testCase.SetUpRequest(ctx, testCase.Body)

			mockServices := testCase.SetupMocksServices(t)

			shortenURLHandler := NewShortenUrlHandler(mockServices)
			shortenURLHandler.CreateShortenUrl(ctx)

			assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.ExpectedResponse, recorder.Body.String())
		})
	}
}

func TestRedirect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		SetupRequest       func(ctx *gin.Context)
		SetupMocksServices func(t *testing.T, ctx *gin.Context) *mocksService.ShortenUrlService
		ExpectedStatusCode int
		ExpectedResponse   string
	}{
		{
			name: "success",
			SetupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/redirect/1234567", nil)
				ctx.Params = gin.Params{
					{
						Key:   "code",
						Value: "1234567",
					},
				}
			},
			SetupMocksServices: func(t *testing.T, ctx *gin.Context) *mocksService.ShortenUrlService {
				mockShortenUrlService := mocksService.NewShortenUrlService(t)
				mockShortenUrlService.On("GetUrlByCode", ctx, "1234567").Return("http://localhost:8000", nil)
				return mockShortenUrlService
			},
			ExpectedStatusCode: http.StatusFound,
			ExpectedResponse:   "http://localhost:8000",
		},
		{
			name: "Param required",
			SetupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/redirect/", nil)
			},
			SetupMocksServices: func(t *testing.T, ctx *gin.Context) *mocksService.ShortenUrlService {
				return mocksService.NewShortenUrlService(t)
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `"Invalid input"`,
		},
		{
			name: "code not found",
			SetupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/redirect/1234567", nil)
				ctx.Params = gin.Params{
					{
						Key:   "code",
						Value: "1234567",
					},
				}
			},
			SetupMocksServices: func(t *testing.T, ctx *gin.Context) *mocksService.ShortenUrlService {
				mockShortenUrlService := mocksService.NewShortenUrlService(t)
				mockShortenUrlService.On("GetUrlByCode", ctx, "1234567").Return("", repository.ErrCodeNotFound)
				return mockShortenUrlService
			},
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   `"Code not found"`,
		},
		{
			name: "internal server error",
			SetupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/redirect/1234567", nil)
				ctx.Params = gin.Params{
					{
						Key:   "code",
						Value: "1234567",
					},
				}
			},
			SetupMocksServices: func(t *testing.T, ctx *gin.Context) *mocksService.ShortenUrlService {
				mockShortenUrlService := mocksService.NewShortenUrlService(t)
				mockShortenUrlService.On("GetUrlByCode", ctx, "1234567").Return("", errors.New("Internal server error"))
				return mockShortenUrlService
			},
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedResponse:   `"Internal server error"`,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			testCase.SetupRequest(ctx)

			mockServices := testCase.SetupMocksServices(t, ctx)

			shortenURLHandler := NewShortenUrlHandler(mockServices)
			shortenURLHandler.Redirect(ctx)

			assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)

			if recorder.Header().Get("Location") != "" {
				assert.Equal(t, testCase.ExpectedResponse, recorder.Header().Get("Location"))
			} else {
				assert.Equal(t, testCase.ExpectedResponse, recorder.Body.String())
			}
		})
	}
}
