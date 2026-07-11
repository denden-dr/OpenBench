package auth

import (
	"fmt"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		if accessToken == "" {
			return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "Access token tidak ditemukan.")
		}

		token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(cfg.Auth.AccessSecret), nil
		})

		if err != nil || !token.Valid {
			return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "Access token tidak valid atau kedaluwarsa.")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "Format token tidak valid.")
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "User ID tidak ditemukan dalam token.")
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
			return utils.SendProblem(c, fiber.StatusForbidden, "/errors/forbidden", "Forbidden Access", "Akses ditolak: role tidak valid.")
		}

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return utils.SendProblem(c, fiber.StatusForbidden, "/errors/forbidden", "Forbidden Access", "Akses ditolak: Anda tidak memiliki akses ke resource ini.")
	}
}
