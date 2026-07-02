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
	return &healthCheckService{
		cfg: cfg,
	}
}

func (h *healthCheckService) CheckHealth() model.HealthCheck {
	instanceId := uuid.New().String()

	if h.cfg.InstanceId != "" {
		instanceId = h.cfg.InstanceId
	}

	healthCheck := model.HealthCheck{
		Message:     "OK",
		ServiceName: h.cfg.ServiceName,
		InstanceId:  instanceId,
	}

	return healthCheck
}
