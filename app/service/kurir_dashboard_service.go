package service

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/repository"
	"time"
)

func GetCourierDashboard(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		courierID := c.Locals("user_id").(int)
		date := time.Now().Format("2006-01-02")

		repo := repository.NewCourierDashboardRepository(db)

		total, pending, delivered, err := repo.GetCourierStats(courierID, date)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"date":             date,
			"total_orders":     total,
			"pending_orders":   pending,
			"delivered_orders": delivered,
		})
	}
}
