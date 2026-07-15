package main

import (
	"github.com/denden-dr/OpenBench/apps/backend/internal/apierrors"
	"github.com/gofiber/fiber/v3"
)

func globalErrorHandler(c fiber.Ctx, err error) error {
	return apierrors.GlobalErrorHandler(c, err)
}
