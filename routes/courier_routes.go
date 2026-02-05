package routes

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/service"
	"masakin-backend/middleware"
)

func RegisterCourierRoutes(app *fiber.App, db *sql.DB) {

	courier := app.Group(
		"api/courier",
		middleware.JWTProtected(),
		middleware.RoleOnly("kurir"),
	)

	courier.Get("/orders", service.GetOrdersForCourier(db))
	courier.Put("/orders/:id/deliver", service.DeliverOrder(db))
	courier.Get("/dashboard", service.GetCourierDashboard(db))
}
