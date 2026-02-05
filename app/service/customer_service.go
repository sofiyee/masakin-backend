package service

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/repository"
	"log"
	
)

func GetCustomerProfile(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userIDRaw := c.Locals("user_id")
		if userIDRaw == nil {
			return fiber.ErrUnauthorized
		}

		userID := userIDRaw.(int)

		repo := repository.NewCustomerRepository(db)
		profile, err := repo.GetProfileByUserID(userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(404).JSON(fiber.Map{
					"error": "customer profile not found",
				})
			}

			log.Println("PROFILE ERROR:", err)
			return fiber.ErrInternalServerError
		}


		return c.JSON(profile)
	}
}

func GetRecentCustomerOrders(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// 1. ambil user_id dari JWT
		userIDRaw := c.Locals("user_id")
		if userIDRaw == nil {
			return fiber.ErrUnauthorized
		}
		userID := userIDRaw.(int)

		// 2. cari customer
		customerRepo := repository.NewCustomerRepository(db)
		customer, err := customerRepo.FindByUserID(userID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "customer not found",
			})
		}

		// 3. ambil recent orders
		orderRepo := repository.NewOrderRepository(db)

		orders, err := orderRepo.GetRecentOrdersByCustomer(customer.ID, 5)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(orders)
	}
}

func GetCustomerOrders(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		customerID := c.Locals("user_id")

		rows, err := db.Query(`
			SELECT 
				id,
				order_type,
				status,
				total_price,
				start_date,
				end_date,
				created_at
			FROM orders
			WHERE customer_id = $1
			ORDER BY created_at DESC
		`, customerID)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer rows.Close()

		var orders []fiber.Map

		for rows.Next() {
			var (
				id          int
				orderType   string
				status      string
				totalPrice  int
				startDate   string
				endDate     string
				createdAt   string
			)

			if err := rows.Scan(
				&id,
				&orderType,
				&status,
				&totalPrice,
				&startDate,
				&endDate,
				&createdAt,
			); err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			orders = append(orders, fiber.Map{
				"order_id":    id,
				"order_type":  orderType,
				"status":      status,
				"total_price": totalPrice,
				"start_date":  startDate,
				"end_date":    endDate,
				"created_at":  createdAt,
			})
		}

		return c.JSON(orders)
	}
}