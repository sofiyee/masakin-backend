package routes

import (
	"database/sql"

	"masakin-backend/app/service"
	"masakin-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterOrderRoutes(app *fiber.App, db *sql.DB) {

	// =========================
	// CUSTOMER ORDER
	// =========================
	customer := app.Group(
		"/api/orders",
		middleware.JWTProtected(),
		middleware.RoleOnly("customer"),
	)

	// create order (daily / monthly)
	customer.Post("/", service.CreateOrder(db))

	// =========================
	// ADMIN ORDER MANAGEMENT
	// =========================
	admin := app.Group(
		"/api/admin/orders",
		middleware.JWTProtected(),
		middleware.RoleOnly("admin"),
	)

	// get all orders
	admin.Get("/", service.GetAllOrders(db))
	admin.Get("/orders/unpaid", service.GetUnpaidOrdersAdmin(db))

	// get order detail
	admin.Get("/:id", service.GetOrderDetail(db))

	

	
}
