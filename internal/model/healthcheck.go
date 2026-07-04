package model

type HealthCheck struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceId  string `json:"instance_id"`
}
