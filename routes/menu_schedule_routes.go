package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/service"
	"masakin-backend/middleware"
)

func RegisterMenuScheduleRoutes(app *fiber.App, db *sql.DB) {

	admin := app.Group(
		"/api/admin/menu-schedules",
		middleware.JWTProtected(),
		middleware.RoleOnly("admin"),
	)

	// ================= ADMIN =================
	admin.Post("/", service.AssignMenuToDate(db))

	admin.Get("/", service.GetMenuSchedulesByMonth(db)) 
	// GET /api/admin/menu-schedules?year=2026&month=2

	admin.Delete("/", service.DeleteMenuSchedule(db))
	// DELETE single menu in date (menu_id + date in body)

	admin.Delete("/date", service.ClearMenuSchedulesByDate(db))
	// DELETE all menus in date (?date=yyyy-mm-dd)

	// ================= CUSTOMER =================
	app.Get("/api/menu-schedule", service.GetMenusByDate(db))
}
