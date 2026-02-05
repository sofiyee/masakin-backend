package routes
import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/service"
	"masakin-backend/middleware"
)

func RegisterCustomerDashboardRoutes(app *fiber.App, db *sql.DB) {
	customer := app.Group(
		"/api/customer",
		middleware.JWTProtected(),
		middleware.RoleOnly("customer"),
	)

	customer.Get("/dashboard", service.GetCustomerDashboard(db))
	customer.Get("/profile", service.GetCustomerProfile(db))
	customer.Get("/orders/recent", service.GetRecentCustomerOrders(db))
	customer.Get("/orders", service.GetCustomerOrders(db))
}
