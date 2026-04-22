package service

import "openbench/server/internal/dto"

type HealthService interface {
	CheckHealth() dto.HealthData
}

type healthServiceImpl struct {
	version string
}

func NewHealthService(version string) HealthService {
	return &healthServiceImpl{
		version: version,
	}
}

func (s *healthServiceImpl) CheckHealth() dto.HealthData {
	return dto.HealthData{
		Version: s.version,
	}
}
