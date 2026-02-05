package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/service"
	"masakin-backend/middleware"

)
func RegisterMenuRoutes(app *fiber.App, db *sql.DB) {

	admin := app.Group(
		"/api/admin/menus",
		middleware.JWTProtected(),
		middleware.RoleOnly("admin"),
	)
	admin.Get("/all", service.GetMenusAdmin(db))
	admin.Post("/", service.CreateMenu(db))
	admin.Put("/:id", service.UpdateMenu(db))
	admin.Delete("/:id", service.DisableMenu(db))
	admin.Put("/:id/enable", service.EnableMenu(db))
	

	// customer
	app.Get("/api/menus", service.GetMenus(db))
}


