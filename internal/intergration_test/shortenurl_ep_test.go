package intergration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/NghiaLeopard/bookmark-management/internal/api"
	"github.com/NghiaLeopard/bookmark-management/internal/handler"
	"github.com/NghiaLeopard/bookmark-management/internal/pkg/redis"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		Body                 handler.ShortenUrlInputBody
		setUpServeHttp       func(t *testing.T, body handler.ShortenUrlInputBody) *httptest.ResponseRecorder
		ExpectedStatusCode   int
		ExpectedResponseBody string
	}{
		{
			name: "success",
			Body: handler.ShortenUrlInputBody{
				Url:    "https://www.google.com",
				Expire: time.Hour,
			},

			setUpServeHttp: func(t *testing.T, body handler.ShortenUrlInputBody) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				bodyJson, err := json.Marshal(body)
				if err != nil {
					assert.NoError(t, err)
				}

				request := httptest.NewRequest(http.MethodPost, "/links/shorten", strings.NewReader(string(bodyJson)))

				mockRedis := redis.NewMockRClient(t)
				app := api.NewEngine(mockRedis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusOK,
			ExpectedResponseBody: `{"code":"`,
		},

		{
			name: "error with empty url",
			Body: handler.ShortenUrlInputBody{
				Expire: time.Hour,
			},

			setUpServeHttp: func(t *testing.T, body handler.ShortenUrlInputBody) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				bodyJson, err := json.Marshal(body)
				if err != nil {
					assert.NoError(t, err)
				}

				request := httptest.NewRequest(http.MethodPost, "/links/shorten", strings.NewReader(string(bodyJson)))

				mockRedis := redis.NewMockRClient(t)
				app := api.NewEngine(mockRedis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusBadRequest,
			ExpectedResponseBody: `Invalid input`,
		},
		{
			name: "error with empty expire",
			Body: handler.ShortenUrlInputBody{
				Url: "https://www.google.com",
			},

			setUpServeHttp: func(t *testing.T, body handler.ShortenUrlInputBody) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				bodyJson, err := json.Marshal(body)
				if err != nil {
					assert.NoError(t, err)
				}

				request := httptest.NewRequest(http.MethodPost, "/links/shorten", strings.NewReader(string(bodyJson)))

				mockRedis := redis.NewMockRClient(t)
				app := api.NewEngine(mockRedis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusBadRequest,
			ExpectedResponseBody: `Invalid input`,
		},

		{
			name: "error with redis is closed",
			Body: handler.ShortenUrlInputBody{
				Url:    "https://www.google.com",
				Expire: time.Hour,
			},

			setUpServeHttp: func(t *testing.T, body handler.ShortenUrlInputBody) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				bodyJson, err := json.Marshal(body)
				if err != nil {
					assert.NoError(t, err)
				}

				request := httptest.NewRequest(http.MethodPost, "/links/shorten", strings.NewReader(string(bodyJson)))

				mockRedis := redis.NewMockRClient(t)
				mockRedis.Close()
				app := api.NewEngine(mockRedis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusInternalServerError,
			ExpectedResponseBody: `"Internal server error"`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			recorder := testCase.setUpServeHttp(t, testCase.Body)

			assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), testCase.ExpectedResponseBody)
		})
	}
}

func TestRedirect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		SetupServerHttp      func(t *testing.T) *httptest.ResponseRecorder
		ExpectedStatusCode   int
		ExpectedResponseBody string
	}{
		{
			name: "success",
			SetupServerHttp: func(t *testing.T) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(http.MethodGet, "/links/redirect/1234567", nil)

				mockRedis := redis.NewMockRClient(t)

				mockRedis.Set(context.Background(), "1234567", "http://localhost:8000", time.Hour)

				app := api.NewEngine(mockRedis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusFound,
			ExpectedResponseBody: "http://localhost:8000",
		},
		{
			name: "error with empty code",
			SetupServerHttp: func(t *testing.T) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(http.MethodGet, "/links/redirect/", nil)

				mockRedis := redis.NewMockRClient(t)

				app := api.NewEngine(mockRedis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusNotFound,
			ExpectedResponseBody: `404 page not found`,
		},
		{
			name: "error with code not found",
			SetupServerHttp: func(t *testing.T) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(http.MethodGet, "/links/redirect/1234567", nil)

				mockRedis := redis.NewMockRClient(t)

				app := api.NewEngine(mockRedis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusNotFound,
			ExpectedResponseBody: `Code not found`,
		},
		{
			name: "error with redis is closed",
			SetupServerHttp: func(t *testing.T) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(http.MethodGet, "/links/redirect/1234567", nil)

				mockRedis := redis.NewMockRClient(t)

				mockRedis.Close()

				app := api.NewEngine(mockRedis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusInternalServerError,
			ExpectedResponseBody: `Internal server error`,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			recorder := testCase.SetupServerHttp(t)

			if recorder.Header().Get("Location") != "" {
				assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
				assert.Equal(t, testCase.ExpectedResponseBody, recorder.Header().Get("Location"))
			} else {
				assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
				assert.Contains(t, recorder.Body.String(), testCase.ExpectedResponseBody)
			}
		})
	}
}
