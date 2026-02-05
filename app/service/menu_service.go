package service

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"masakin-backend/app/model"
	"masakin-backend/app/repository"

	"github.com/gofiber/fiber/v2"
)

// CREATE MENU (ADMIN)
func CreateMenu(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		name := c.FormValue("name")
		desc := c.FormValue("description")
		price, _ := strconv.Atoi(c.FormValue("price"))
		month, _ := strconv.Atoi(c.FormValue("menu_month"))
		year, _ := strconv.Atoi(c.FormValue("menu_year"))

		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "image required"})
		}

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		path := filepath.Join("uploads/menus", filename)

		if err := c.SaveFile(file, path); err != nil {
			return fiber.ErrInternalServerError
		}

		repo := repository.NewMenuRepository(db)
		menu := &model.Menu{
			Name:        name,
			Description: desc,
			Price:       price,
			ImageURL:    path,
			MenuMonth:   month,
			MenuYear:    year,
		}

		if err := repo.Create(menu); err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"message": "menu created",
		})
	}
}


// get menu
func GetMenus(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		month, _ := strconv.Atoi(c.Query("month"))
		year, _ := strconv.Atoi(c.Query("year"))

		if month == 0 || year == 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "month and year required",
			})
		}

		repo := repository.NewMenuRepository(db)
		data, err := repo.GetByMonthYear(month, year)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(data)
	}
}
// UPDATE MENU
func UpdateMenu(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Price       int    `json:"price"`
			ImageURL    string `json:"image_url"`
			MenuMonth   int    `json:"menu_month"`
			MenuYear    int    `json:"menu_year"`
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		// VALIDASI WAJIB
		if req.MenuMonth == 0 || req.MenuYear == 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "menu_month and menu_year required",
			})
		}

		repo := repository.NewMenuRepository(db)

		err = repo.Update(id, &model.Menu{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			ImageURL:    req.ImageURL,
			MenuMonth:   req.MenuMonth,
			MenuYear:    req.MenuYear,
		})
		if err != nil {
			// 🔥 INI PENTING BUAT DEBUG
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "menu updated",
		})
	}
}

// DISABLE MENU
func DisableMenu(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		repo := repository.NewMenuRepository(db)
		if err := repo.Disable(id); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "menu disabled",
		})
	}
}
// ENABLE MENU

func EnableMenu(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		repo := repository.NewMenuRepository(db)
		if err := repo.Enable(id); err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"message": "menu enabled",
		})
	}
}

// GET MENUS FOR ADMIN
func GetMenusAdmin(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		month, _ := strconv.Atoi(c.Query("month"))
		year, _ := strconv.Atoi(c.Query("year"))

		if month == 0 || year == 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "month and year required",
			})
		}

		repo := repository.NewMenuRepository(db)
		data, err := repo.GetByMonthYearAdmin(month, year)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(data)
	}
}