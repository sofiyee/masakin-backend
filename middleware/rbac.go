package middleware

import "github.com/gofiber/fiber/v2"


func RoleOnly(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		roleVal := c.Locals("role")
		if roleVal == nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		

		role := roleVal.(string)


		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "access denied",
		})
	}
}
