package service

import (
	"errors"
	"testing"

	"github.com/NghiaLeopard/bookmark-management/internal/config"
	"github.com/NghiaLeopard/bookmark-management/internal/model"
	redisPkg "github.com/NghiaLeopard/bookmark-management/internal/pkg/redis"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                    string
		SetupMockRedis          func(t *testing.T) *redis.Client
		ExpectedConfig          *config.Config
		ExpectedBodyHealthCheck model.HealthCheck
		ExpectedMessageResponse string
		ExpectedRedisStatus     string
		ExpectedError           error
	}{
		{
			name: "Success",
			SetupMockRedis: func(t *testing.T) *redis.Client {
				mockRedis := redisPkg.NewMockRClient(t)
				return mockRedis
			},
			ExpectedConfig:          &config.Config{Port: "8080", ServiceName: "bookmark-management", InstanceId: "1234567890"},
			ExpectedBodyHealthCheck: model.HealthCheck{Message: "OK", ServiceName: "bookmark-management", InstanceId: "1234567890", RedisStatus: "PONG"},
			ExpectedMessageResponse: "OK",
			ExpectedRedisStatus:     "PONG",
			ExpectedError:           nil,
		},
		{
			name: "Success with empty service name",
			SetupMockRedis: func(t *testing.T) *redis.Client {
				mockRedis := redisPkg.NewMockRClient(t)
				return mockRedis
			},
			ExpectedConfig:          &config.Config{Port: "8080", ServiceName: "", InstanceId: "1234567890"},
			ExpectedBodyHealthCheck: model.HealthCheck{Message: "OK", ServiceName: "", InstanceId: "1234567890", RedisStatus: "PONG"},
			ExpectedMessageResponse: "OK",
			ExpectedRedisStatus:     "PONG",
			ExpectedError:           nil,
		},
		{
			name: "Success with empty instance id",
			SetupMockRedis: func(t *testing.T) *redis.Client {
				mockRedis := redisPkg.NewMockRClient(t)
				return mockRedis
			},
			ExpectedConfig:          &config.Config{Port: "8080", ServiceName: "bookmark-management", InstanceId: ""},
			ExpectedBodyHealthCheck: model.HealthCheck{Message: "OK", ServiceName: "bookmark-management", InstanceId: uuid.New().String(), RedisStatus: "PONG"},
			ExpectedMessageResponse: "OK",
			ExpectedRedisStatus:     "PONG",
			ExpectedError:           nil,
		},
		{
			name: "Success with empty instance id and service name",
			SetupMockRedis: func(t *testing.T) *redis.Client {
				mockRedis := redisPkg.NewMockRClient(t)
				return mockRedis
			},
			ExpectedConfig:          &config.Config{Port: "8080", ServiceName: "", InstanceId: ""},
			ExpectedBodyHealthCheck: model.HealthCheck{Message: "OK", ServiceName: "", InstanceId: uuid.New().String(), RedisStatus: "PONG"},
			ExpectedMessageResponse: "OK",
			ExpectedRedisStatus:     "PONG",
			ExpectedError:           nil,
		},
		{
			name: "Error with redis is closed",
			SetupMockRedis: func(t *testing.T) *redis.Client {
				mockRedis := redisPkg.NewMockRClient(t)

				mockRedis.Close()
				return mockRedis
			},
			ExpectedConfig:          &config.Config{Port: "8080", ServiceName: "bookmark-management", InstanceId: "1234567890"},
			ExpectedBodyHealthCheck: model.HealthCheck{Message: "Error", ServiceName: "bookmark-management", InstanceId: "1234567890", RedisStatus: "Error"},
			ExpectedMessageResponse: "Error",
			ExpectedRedisStatus:     "Error",
			ExpectedError:           errors.New("redis: client is closed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockRedis := testCase.SetupMockRedis(t)
			testService := NewHealthCheck(testCase.ExpectedConfig, mockRedis)
			healthCheck, err := testService.CheckHealth()

			assert.Equal(t, testCase.ExpectedError, err)
			assert.Equal(t, testCase.ExpectedMessageResponse, healthCheck.Message)
			assert.Equal(t, testCase.ExpectedConfig.ServiceName, healthCheck.ServiceName)
			assert.NotEmpty(t, healthCheck.InstanceId)
			assert.Equal(t, len(testCase.ExpectedBodyHealthCheck.InstanceId), len(healthCheck.InstanceId))
			assert.Equal(t, testCase.ExpectedRedisStatus, healthCheck.RedisStatus)
		})
	}
}
