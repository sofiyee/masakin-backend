package middleware

import (
	"strings"

	"masakin-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "missing token",
			})
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token format",
			})
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		
		c.Locals("user_id", claims.UserID) // INT
		c.Locals("name", claims.Name)      // STRING
		c.Locals("role", claims.Role)      // STRING

		return c.Next()
	}
}
