package auth

import (
	"fmt"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(cfg *config.Config, queryRepo QueryRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		if accessToken == "" {
			return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "Access token is missing.")
		}

		token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(cfg.Auth.AccessSecret), nil
		})

		if err != nil || !token.Valid {
			return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "Access token is invalid or expired.")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "Invalid token format.")
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "User ID not found in token.")
		}

		jti, ok := claims["jti"].(string)
		if ok && jti != "" {
			isBlacklisted, err := queryRepo.IsTokenBlacklisted(c.Context(), jti)
			if err != nil || isBlacklisted {
				return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "Access token has been revoked.")
			}
		}

		role, _ := claims["role"].(string)

		c.Locals("userID", userID)
		c.Locals("userRole", role)

		return c.Next()
	}
}

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRole, ok := c.Locals("userRole").(string)
		if !ok || userRole == "" {
			return utils.SendProblem(c, fiber.StatusForbidden, "/errors/forbidden", "Forbidden Access", "Access denied: invalid role.")
		}

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return utils.SendProblem(c, fiber.StatusForbidden, "/errors/forbidden", "Forbidden Access", "Access denied: you do not have permission to access this resource.")
	}
}
