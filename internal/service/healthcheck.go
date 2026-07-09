package service

import (
	"context"

	"github.com/NghiaLeopard/bookmark-management/internal/config"
	"github.com/NghiaLeopard/bookmark-management/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name HealthCheck --filename healthcheck.go
type HealthCheck interface {
	CheckHealth() (model.HealthCheck, error)
}

type healthCheckService struct {
	cfg *config.Config
	rdb *redis.Client
}

func NewHealthCheck(cfg *config.Config, rdb *redis.Client) HealthCheck {
	if cfg.InstanceId == "" {
		instanceId := uuid.New().String()

		cfg.InstanceId = instanceId
	}

	return &healthCheckService{
		cfg: cfg,
		rdb: rdb,
	}
}

func (h *healthCheckService) CheckHealth() (model.HealthCheck, error) {
	result, err := h.rdb.Ping(context.Background()).Result()

	if err != nil || result != "PONG" {
		return model.HealthCheck{
			Message:     "Error",
			ServiceName: h.cfg.ServiceName,
			InstanceId:  h.cfg.InstanceId,
			RedisStatus: "Error",
		}, err
	}

	healthCheck := model.HealthCheck{
		Message:     "OK",
		ServiceName: h.cfg.ServiceName,
		InstanceId:  h.cfg.InstanceId,
		RedisStatus: "PONG",
	}

	return healthCheck, nil
}
