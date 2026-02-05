package service
import (
	"database/sql"
	"masakin-backend/app/repository"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/model"
	"masakin-backend/utils"
	"strconv"
	
)
// CREATE ORDER CUSTOMER
func CreateOrder(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// =====================
		// 1. AMBIL USER_ID DARI JWT
		// =====================
		userIDRaw := c.Locals("user_id")
		if userIDRaw == nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		userID, ok := userIDRaw.(int)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token payload",
			})
		}

		// =====================
		// 2. CARI CUSTOMER
		// =====================
		customerRepo := repository.NewCustomerRepository(db)
		customer, err := customerRepo.FindByUserID(userID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "customer profile not found",
			})
		}

		// =====================
		// 3. PARSE REQUEST
		// =====================
		var req struct {
			OrderType string `json:"order_type"` // daily | monthly
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
			Items []struct {
				MenuID int    `json:"menu_id"`
				Date   string `json:"order_date"`
				Qty    int    `json:"quantity"`
			} `json:"items"`
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		if len(req.Items) == 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "order items required",
			})
		}

		startDate, err := utils.ParseDate(req.StartDate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid start_date",
			})
		}

		endDate, err := utils.ParseDate(req.EndDate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid end_date",
			})
		}

		// =====================
		// 4. TRANSACTION START
		// =====================
		tx, err := db.Begin()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		defer tx.Rollback()

		orderRepo := repository.NewOrderRepository(tx)
		itemRepo  := repository.NewOrderItemRepository(tx)
		menuRepo  := repository.NewMenuRepository(tx)


		// =====================
		// 5. CREATE ORDER
		// =====================
		
		orderID, err := orderRepo.Create(&model.Order{
			CustomerID: customer.ID,
			OrderType:  req.OrderType,
			StartDate:  startDate,
			EndDate:    endDate,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		total := 0

		// =====================
		// 6. INSERT ITEMS
		// =====================
		for _, item := range req.Items {

			if item.Qty <= 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "quantity must be greater than 0",
				})
			}

			price, err := menuRepo.GetPriceByID(item.MenuID)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{
					"error": "invalid menu_id",
				})
			}

			orderDate, err := utils.ParseDate(item.Date)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{
					"error": "invalid order_date",
				})
			}

			err = itemRepo.Create(&model.OrderItem{
				OrderID:   orderID,
				MenuID:    item.MenuID,
				OrderDate: orderDate,
				Quantity:  item.Qty,
			})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			total += price * item.Qty
		}

		// =====================
		// 7. UPDATE TOTAL
		// =====================
		if err := orderRepo.UpdateTotal(orderID, total); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// =====================
		// 8. COMMIT
		// =====================
		if err := tx.Commit(); err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"message":  "order created",
			"order_id": orderID,
			"total":    total,
		})
	}
}

// GET ALL ORDERS (ADMIN)
func GetAllOrders(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		repo := repository.NewOrderRepository(db)

		orders, err := repo.GetAll()
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(orders)
	}
}

// Get Order Details (Admin)
func GetOrderDetail(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		orderID, _ := strconv.Atoi(c.Params("id"))

		orderRepo := repository.NewOrderRepository(db)
		itemRepo := repository.NewOrderItemRepository(db)

		order, err := orderRepo.GetByID(orderID)
		if err != nil {
			return fiber.ErrNotFound
		}

		items, err := itemRepo.GetDetail(orderID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"order": order,
			"items": items,
		})
	}
}

// =====================
// GET RECENT ORDERS (Admin Dashboard)
// =====================
func GetRecentOrders(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get limit from query parameter, default to 5
		limitStr := c.Query("limit", "5")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 5
		}
		
		// Max limit to prevent overload
		if limit > 50 {
			limit = 50
		}

		repo := repository.NewOrderRepository(db)
		orders, err := repo.GetRecentOrders(limit)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to fetch recent orders",
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    orders,
			"message": "Recent orders fetched successfully",
		})
	}
}

func GetUnpaidOrdersAdmin(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		repo := repository.NewOrderRepository(db)

		orders, err := repo.GetUnpaidOrdersAdmin()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(orders)
	}
}

