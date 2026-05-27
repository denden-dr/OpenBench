package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/denden-dr/openbench/apps/backend/internal/dto"
	"github.com/denden-dr/openbench/apps/backend/internal/handler"
	"github.com/denden-dr/openbench/apps/backend/internal/middleware"
	mockservice "github.com/denden-dr/openbench/apps/backend/mocks/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTicketHandler_GetByID(t *testing.T) {
	ticketID := "11111111-1111-1111-1111-111111111111"

	t.Run("success", func(t *testing.T) {
		mockService := new(mockservice.MockTicketService)
		expectedResponse := &dto.TicketResponse{
			ID:           ticketID,
			CustomerName: "Budi",
			Brand:        "Apple",
			Model:        "iPhone 13",
		}
		mockService.On("GetTicket", mock.Anything, ticketID).Return(expectedResponse, nil).Once()

		h := handler.NewTicketHandler(mockService)
		app := fiber.New(fiber.Config{
			ErrorHandler: middleware.ErrorHandler,
		})
		app.Get("/tickets/:id", h.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/tickets/"+ticketID, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.True(t, result["success"].(bool))

		data := result["data"].(map[string]interface{})
		assert.Equal(t, ticketID, data["id"])
		assert.Equal(t, "Budi", data["customer_name"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid UUID format", func(t *testing.T) {
		mockService := new(mockservice.MockTicketService)
		h := handler.NewTicketHandler(mockService)
		app := fiber.New(fiber.Config{
			ErrorHandler: middleware.ErrorHandler,
		})
		app.Get("/tickets/:id", h.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/tickets/invalid-id-format", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}
