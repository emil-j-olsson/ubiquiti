package service

import "github.com/emil-j-olsson/ubiquiti/device/internal/types"

// retrieve data from cache... (update cache via env variables + api calls...)
type StateRetriever interface {
	GetState() types.DeviceState
}

type Service struct {
}

// Do fundamental service
func NewDeviceService() *Service {
	return &Service{}
}

func (s *Service) GetHealth() *types.HealthStatus {
	// implement core business logic here
	// grab settings from cache (memorystore)
	return &types.HealthStatus{}
}

// get checksum from internal module of external module (binary)
