package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/service"
	"masakin-backend/middleware"

)

func RegisterAdminRoutes(app *fiber.App, db *sql.DB) {

	admin := app.Group(
		"/api/admin",
		middleware.JWTProtected(),
		middleware.RoleOnly("admin"),
	)

	admin.Post("/couriers", service.CreateCourier(db))
	admin.Get("/couriers", service.GetCouriers(db))
	admin.Put("/couriers/region", service.UpdateCourierRegion(db))
	admin.Put("/couriers/status", service.SetCourierActive(db))
	admin.Post("/gudang", service.CreateGudangUser(db))
	admin.Get("/reports/monthly", service.GetMonthlyReport(db))
	admin.Get("/dashboard", service.GetAdminDashboard(db))
	admin.Get("/dashboard/recent-orders", service.GetRecentOrders(db))
}
