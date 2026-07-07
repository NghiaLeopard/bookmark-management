package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
		ExpectedError      error
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
			ExpectedError:      nil,
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
			ExpectedError:      errors.New("Invalid input"),
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
			ExpectedError:      errors.New("Invalid input"),
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
			ExpectedError:      errors.New("Internal server error"),
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
