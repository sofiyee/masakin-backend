package service

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/repository"
	"time"
)

func GetAdminDashboard(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		month := time.Now().Format("2006-01")
		repo := repository.NewDashboardRepository(db)

		todayOrders, err := repo.GetTodayOrders()
		if err != nil {
			return fiber.ErrInternalServerError
		}

		monthlyOrders, err := repo.GetMonthlyOrders(month)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		monthlyRevenue, err := repo.GetMonthlyRevenue(month)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		activeCustomer, err := repo.GetActiveCustomer()
		if err != nil {
			return fiber.ErrInternalServerError
		}

		activeCourier, err := repo.GetActiveCourier()
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"today_orders":     todayOrders,
			"monthly_orders":   monthlyOrders,
			"monthly_revenue":  monthlyRevenue,
			"active_customer":  activeCustomer,
			"active_courier":   activeCourier,
		})
	}
}
