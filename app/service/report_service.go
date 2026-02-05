package service

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/repository"
	"log"
)

func GetMonthlyReport(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		month := c.Query("month")
		if month == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "month is required (YYYY-MM)",
			})
		}

		repo := repository.NewReportRepository(db)

		totalPortion, err := repo.GetTotalPortion(month)
		if err != nil {
			log.Println("ERROR GetTotalPortion:", err)
			return fiber.ErrInternalServerError
		}

		totalRevenue, err := repo.GetTotalRevenue(month)
		if err != nil {
			log.Println("ERROR GetTotalRevenue:", err)
			return fiber.ErrInternalServerError
		}

		topMenus, err := repo.GetTopMenus(month)
		if err != nil {
			log.Println("ERROR GetTopMenus:", err)
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"month":          month,
			"total_portion":  totalPortion,
			"total_revenue":  totalRevenue,
			"top_menus":      topMenus,
		})
	}
}


