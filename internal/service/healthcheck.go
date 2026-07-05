package service

import (
	"github.com/NghiaLeopard/bookmark-management/internal/config"
	"github.com/NghiaLeopard/bookmark-management/internal/model"
	"github.com/google/uuid"
)

//go:generate mockery --name HealthCheck --filename healthcheck.go
type HealthCheck interface {
	CheckHealth() model.HealthCheck
}

type healthCheckService struct {
	cfg *config.Config
}

func NewHealthCheck(cfg *config.Config) HealthCheck {
	if cfg.InstanceId == "" {
		instanceId := uuid.New().String()

		cfg.InstanceId = instanceId
	}

	return &healthCheckService{
		cfg: cfg,
	}
}

func (h *healthCheckService) CheckHealth() model.HealthCheck {
	healthCheck := model.HealthCheck{
		Message:     "OK",
		ServiceName: h.cfg.ServiceName,
		InstanceId:  h.cfg.InstanceId,
	}

	return healthCheck
}
