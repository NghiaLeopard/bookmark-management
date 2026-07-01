package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NghiaLeopard/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testError = errors.New("Internal server error")

func TestGeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setUpRequest     func(ctx *gin.Context)
		setUpMockService func(ctx context.Context) *mocks.GenPass

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Success",
			setUpRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/genpass", nil)
			},
			setUpMockService: func(ctx context.Context) *mocks.GenPass {
				serviceMock := mocks.NewGenPass(t)
				serviceMock.On("GeneratePassword", 12).Return("123456789012", nil)

				return serviceMock
			},

			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"password":"123456789012"}`,
		},
		{
			name: "Service fail",
			setUpRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/genpass", nil)
			},
			setUpMockService: func(ctx context.Context) *mocks.GenPass {
				serviceMock := mocks.NewGenPass(t)
				serviceMock.On("GeneratePassword", 12).Return("", testError)

				return serviceMock
			},

			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal server error"}`,
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

			testHandler := NewGenPassHandler(serviceMock)
			testHandler.GeneratePassword(ctx)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponse, recorder.Body.String())
		})
	}
}
