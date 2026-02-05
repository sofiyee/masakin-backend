package service

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/repository"
	"time"
)

func GetWarehouseDashboard(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		date := time.Now().Format("2006-01-02")
		repo := repository.NewWarehouseDashboardRepository(db)

		totalPortion, err := repo.GetTodayTotalPortion(date)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		totalMenu, err := repo.GetTodayTotalMenu(date)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"date":          date,
			"total_portion": totalPortion,
			"total_menu":    totalMenu,
		})
	}
}
