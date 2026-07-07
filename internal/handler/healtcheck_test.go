package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NghiaLeopard/bookmark-management/internal/model"
	"github.com/NghiaLeopard/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCheckHealth(t *testing.T) {
	t.Parallel()
	instanceId := uuid.New().String()

	testCases := []struct {
		name             string
		setUpRequest     func(ctx *gin.Context)
		setUpMockService func(ctx context.Context) *mocks.HealthCheck

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Success",
			setUpRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setUpMockService: func(ctx context.Context) *mocks.HealthCheck {
				serviceMock := mocks.NewHealthCheck(t)
				serviceMock.On("CheckHealth").Return(model.HealthCheck{Message: "OK", ServiceName: "bookmark-management", InstanceId: "1234567890", RedisStatus: "PONG"}, nil)

				return serviceMock
			},

			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"OK","service_name":"bookmark-management","instance_id":"1234567890","redis_status":"PONG"}`,
		},
		{
			name: "Success with empty service name",
			setUpRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},

			setUpMockService: func(ctx context.Context) *mocks.HealthCheck {
				serviceMock := mocks.NewHealthCheck(t)
				serviceMock.On("CheckHealth").Return(model.HealthCheck{Message: "OK", ServiceName: "", InstanceId: "1234567890", RedisStatus: "PONG"}, nil)

				return serviceMock
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"OK","service_name":"","instance_id":"1234567890","redis_status":"PONG"}`,
		},
		{
			name: "Success with empty instance id",
			setUpRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},

			setUpMockService: func(ctx context.Context) *mocks.HealthCheck {
				serviceMock := mocks.NewHealthCheck(t)
				serviceMock.On("CheckHealth").Return(model.HealthCheck{Message: "OK", ServiceName: "bookmark-management", InstanceId: instanceId, RedisStatus: "PONG"}, nil)

				return serviceMock
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`{"message":"OK","service_name":"bookmark-management","instance_id":"%s","redis_status":"PONG"}`, instanceId),
		},
		{
			name: "Success with empty service name and instance id",
			setUpRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},

			setUpMockService: func(ctx context.Context) *mocks.HealthCheck {
				serviceMock := mocks.NewHealthCheck(t)
				serviceMock.On("CheckHealth").Return(model.HealthCheck{Message: "OK", ServiceName: "", InstanceId: instanceId, RedisStatus: "PONG"}, nil)

				return serviceMock
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`{"message":"OK","service_name":"","instance_id":"%s","redis_status":"PONG"}`, instanceId),
		},

		{
			name: "Error with redis not responding",
			setUpRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},

			setUpMockService: func(ctx context.Context) *mocks.HealthCheck {
				serviceMock := mocks.NewHealthCheck(t)
				serviceMock.On("CheckHealth").Return(model.HealthCheck{Message: "Error", ServiceName: "", InstanceId: instanceId, RedisStatus: "Error"}, errors.New("redis: client is closed"))

				return serviceMock
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `"Internal server error"`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Nhận giá trị từ response writer
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			testCase.setUpRequest(ctx)

			serviceMock := testCase.setUpMockService(ctx)

			testHandler := NewHealthCheck(serviceMock)
			testHandler.CheckHealth(ctx)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponse, recorder.Body.String())
		})
	}
}
