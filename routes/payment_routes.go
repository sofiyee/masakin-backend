package routes

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/service"
	"masakin-backend/middleware"
)

func RegisterPaymentRoutes(app *fiber.App, db *sql.DB) {
	admin := app.Group(
		"/api/admin",
		middleware.JWTProtected(),
		middleware.RoleOnly("admin"),
	)
	admin.Get("/payments", service.GetAllPaymentsAdmin(db))
	admin.Post("/payments", service.CreatePayment(db))
	admin.Put("/payments/image", service.UpdatePaymentImage(db))
	

}
