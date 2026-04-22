package service

import (
	"openbench/server/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockHealthService struct {
	mock.Mock
}

func (m *MockHealthService) CheckHealth() dto.HealthData {
	args := m.Called()
	return args.Get(0).(dto.HealthData)
}
