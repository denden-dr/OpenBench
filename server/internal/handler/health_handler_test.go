package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"openbench/server/internal/dto"
	"openbench/server/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestGetHealth(t *testing.T) {
	app := fiber.New()
	mockService := new(service.MockHealthService)
	healthHandler := NewHealthHandler(mockService)

	app.Get("/api/health", healthHandler.GetHealth)

	expectedData := dto.HealthData{Version: "0.1.0"}
	mockService.On("CheckHealth").Return(expectedData)

	req := httptest.NewRequest("GET", "/api/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "OK", response.Message)
	
	// Convert map to struct for comparison if needed, or check fields
	dataMap := response.Data.(map[string]interface{})
	assert.Equal(t, "0.1.0", dataMap["version"])

	mockService.AssertExpectations(t)
}
