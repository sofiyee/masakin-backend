package service

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/repository"
	"time"
)

func GetCustomerDashboard(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userID := c.Locals("user_id").(int)

		customerRepo := repository.NewCustomerRepository(db)
		customer, err := customerRepo.FindByUserID(userID)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		repo := repository.NewCustomerDashboardRepository(db)
		month := time.Now().Format("2006-01")

		hasActive, nextDate, err := repo.GetActiveOrder(customer.ID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		monthlyOrders, err := repo.GetMonthlyOrders(customer.ID, month)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		monthlySpending, err := repo.GetMonthlySpending(customer.ID, month)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"has_active_order":   hasActive,
			"next_delivery_date": nextDate,
			"monthly_orders":     monthlyOrders,
			"monthly_spending":   monthlySpending,
		})
	}
}

