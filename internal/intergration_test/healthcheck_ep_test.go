package intergration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/NghiaLeopard/bookmark-management/internal/api"
	"github.com/NghiaLeopard/bookmark-management/internal/model"
	redisPkg "github.com/NghiaLeopard/bookmark-management/internal/pkg/redis"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckEP(t *testing.T) {
	var instanceId = uuid.New().String()

	testCases := []struct {
		name           string
		setUpEnv       func()
		setUpServeHttp func(t *testing.T) *httptest.ResponseRecorder

		ExpectedStatusCode   int
		ExpectedResponseBody string
	}{
		{
			name: "success",
			setUpEnv: func() {
				t.Setenv("SERVICE_NAME", "bookmark-management")
				t.Setenv("INSTANCE_ID", "1234567890")
			},

			setUpServeHttp: func(t *testing.T) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)

				redis := redisPkg.NewMockRClient(t)
				app := api.NewEngine(redis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusOK,
			ExpectedResponseBody: `{"message":"OK","service_name":"bookmark-management","instance_id":"1234567890","redis_status":"PONG"}`,
		},
		{
			name: "success with empty service name",
			setUpEnv: func() {
				t.Setenv("SERVICE_NAME", "")
				t.Setenv("INSTANCE_ID", instanceId)
			},

			setUpServeHttp: func(t *testing.T) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)

				redis := redisPkg.NewMockRClient(t)
				app := api.NewEngine(redis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusOK,
			ExpectedResponseBody: `{"message":"OK","service_name":"","instance_id":"1234567890","redis_status":"PONG"}`,
		},

		{
			name: "success with empty instance id",
			setUpEnv: func() {
				t.Setenv("SERVICE_NAME", "bookmark-management")
				t.Setenv("INSTANCE_ID", "")
			},

			setUpServeHttp: func(t *testing.T) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)

				redis := redisPkg.NewMockRClient(t)
				app := api.NewEngine(redis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusOK,
			ExpectedResponseBody: `{"message":"OK","service_name":"bookmark-management","instance_id": "` + instanceId + `","redis_status":"PONG"}`,
		},

		{
			name: "success with empty instance id and service name",
			setUpEnv: func() {
				t.Setenv("SERVICE_NAME", "")
				t.Setenv("INSTANCE_ID", "")
			},

			setUpServeHttp: func(t *testing.T) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)

				redis := redisPkg.NewMockRClient(t)
				app := api.NewEngine(redis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusOK,
			ExpectedResponseBody: `{"message":"OK","service_name":"","instance_id":"` + instanceId + `","redis_status":"PONG"}`,
		},

		{
			name: "error with redis is closed",
			setUpEnv: func() {
				t.Setenv("SERVICE_NAME", "bookmark-management")
				t.Setenv("INSTANCE_ID", instanceId)
			},

			setUpServeHttp: func(t *testing.T) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)

				redis := redisPkg.NewMockRClient(t)
				redis.Close()
				app := api.NewEngine(redis)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusInternalServerError,
			ExpectedResponseBody: `"Internal server error"`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setUpEnv()
			recorder := testCase.setUpServeHttp(t)

			if testCase.ExpectedStatusCode == http.StatusOK {
				var marshalData model.HealthCheck
				assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
				assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &marshalData))
				assert.Equal(t, "OK", marshalData.Message)

				assert.Equal(t, os.Getenv("SERVICE_NAME"), marshalData.ServiceName)
				assert.Equal(t, "PONG", marshalData.RedisStatus)
				if os.Getenv("INSTANCE_ID") != "" {
					assert.Equal(t, os.Getenv("INSTANCE_ID"), marshalData.InstanceId)
				} else {
					assert.NotEmpty(t, marshalData.InstanceId)
					assert.Equal(t, 36, len(marshalData.InstanceId))
				}
			} else {
				assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
				assert.Equal(t, testCase.ExpectedResponseBody, recorder.Body.String())
			}
		})
	}
}
