package utils

import (
	"github.com/gofiber/fiber/v3"
)

type ProblemDetail struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

func SendProblem(c fiber.Ctx, status int, problemType, title, detail string) error {
	c.Set("Content-Type", "application/problem+json")
	return c.Status(status).JSON(ProblemDetail{
		Type:     problemType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: c.Path(),
	})
}
