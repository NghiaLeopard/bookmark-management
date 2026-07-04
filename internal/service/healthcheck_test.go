package service

import (
	"testing"

	"github.com/NghiaLeopard/bookmark-management/internal/config"
	"github.com/NghiaLeopard/bookmark-management/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                    string
		ExpectedConfig          *config.Config
		ExpectedBodyHealthCheck model.HealthCheck
	}{
		{
			name:                    "Success",
			ExpectedConfig:          &config.Config{Port: "8080", ServiceName: "bookmark-management", InstanceId: "1234567890"},
			ExpectedBodyHealthCheck: model.HealthCheck{Message: "OK", ServiceName: "bookmark-management", InstanceId: "1234567890"},
		},
		{
			name:                    "Success with empty service name",
			ExpectedConfig:          &config.Config{Port: "8080", ServiceName: "", InstanceId: "1234567890"},
			ExpectedBodyHealthCheck: model.HealthCheck{Message: "OK", ServiceName: "", InstanceId: "1234567890"},
		},
		{
			name:                    "Success with empty instance id",
			ExpectedConfig:          &config.Config{Port: "8080", ServiceName: "bookmark-management", InstanceId: ""},
			ExpectedBodyHealthCheck: model.HealthCheck{Message: "OK", ServiceName: "bookmark-management", InstanceId: uuid.New().String()},
		},
		{
			name:                    "Success with empty instance id and service name",
			ExpectedConfig:          &config.Config{Port: "8080", ServiceName: "", InstanceId: ""},
			ExpectedBodyHealthCheck: model.HealthCheck{Message: "OK", ServiceName: "", InstanceId: uuid.New().String()},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testService := NewHealthCheck(testCase.ExpectedConfig)
			healthCheck := testService.CheckHealth()

			assert.Equal(t, "OK", healthCheck.Message)
			assert.Equal(t, testCase.ExpectedConfig.ServiceName, healthCheck.ServiceName)
			assert.NotEmpty(t, healthCheck.InstanceId)
			assert.Equal(t, len(testCase.ExpectedBodyHealthCheck.InstanceId), len(healthCheck.InstanceId))
		})
	}
}
