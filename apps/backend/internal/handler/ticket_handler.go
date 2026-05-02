package handler

import (
    "github.com/denden-dr/openbench/apps/backend/internal/dto"
    "github.com/denden-dr/openbench/apps/backend/internal/repository"
    "github.com/denden-dr/openbench/apps/backend/internal/service"
    "github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
    service service.TicketService
}

func NewTicketHandler(service service.TicketService) *TicketHandler {
    return &TicketHandler{service: service}
}

func (h *TicketHandler) Create(c *fiber.Ctx) error {
    var req dto.CreateTicketRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "error":   "Invalid request body",
        })
    }

    res, err := h.service.CreateTicket(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "data":    res,
    })
}

func (h *TicketHandler) GetByID(c *fiber.Ctx) error {
    id := c.Params("id")
    res, err := h.service.GetTicket(c.Context(), id)
    if err != nil {
        if err == repository.ErrNotFound {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "success": false,
                "error":   "Ticket not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "error":   "Internal server error",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    res,
    })
}
