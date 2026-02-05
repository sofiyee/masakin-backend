package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/service"
	"masakin-backend/middleware"
	
)

func Register(app *fiber.App, db *sql.DB) {

	api := app.Group("/api")

	// health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
		})
	})


	// AUTH ROUTES
	api.Post("/login", service.Login(db))
	api.Post("/register", service.RegisterCustomer(db))
	api.Post("/logout", middleware.JWTProtected(), service.Logout())
	



	// other route groups
	RegisterAdminRoutes(app, db)
	RegisterMenuRoutes(app, db)
	RegisterMenuScheduleRoutes(app, db)
	RegisterOrderRoutes(app, db)
	RegisterPaymentRoutes(app, db)
	RegisterCourierRoutes(app, db)
	RegisterKitchenRoutes(app, db)
	RegisterCustomerDashboardRoutes(app, db)
	RegisterCaptchaRoutes(app)

}
