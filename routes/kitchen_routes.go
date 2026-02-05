package routes
import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/service"
	"masakin-backend/middleware"
)
func RegisterKitchenRoutes(app *fiber.App, db *sql.DB) {
	warehouse := app.Group(
		"/api/warehouse",
		middleware.JWTProtected(),
		middleware.RoleOnly("gudang", "admin"),
	)

	warehouse.Get("/daily", service.GetDailyKitchenSummary(db))
	warehouse.Get("/dashboard", service.GetWarehouseDashboard(db))

}
