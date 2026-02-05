package service

import (
	"database/sql"
	"log"
	"strconv"

	"masakin-backend/app/repository"

	"github.com/gofiber/fiber/v2"
)


// ================================
// ADMIN - assign menu ke tanggal
// ================================
func AssignMenuToDate(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req struct {
			MenuID int    `json:"menu_id"`
			Date   string `json:"date"` // yyyy-mm-dd
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		repo := repository.NewMenuScheduleRepository(db)

		if err := repo.Create(req.MenuID, req.Date); err != nil {
			return c.Status(409).JSON(fiber.Map{
				"error": "menu already assigned for this date",
			})
		}

		log.Printf("📅 Menu %d assigned to %s", req.MenuID, req.Date)

		return c.JSON(fiber.Map{
			"message": "menu assigned",
		})
	}
}


// ================================
// CUSTOMER - get menu by date
// ================================
func GetMenusByDate(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		date := c.Query("date")
		if date == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "date required",
			})
		}

		repo := repository.NewMenuScheduleRepository(db)
		menus, err := repo.GetByDate(date)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		if len(menus) == 0 {
			return c.Status(404).JSON(fiber.Map{
				"error": "no menu for this date",
			})
		}

		return c.JSON(menus)
	}
}


// ================================
// ADMIN - get schedules by month
// ================================
func GetMenuSchedulesByMonth(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		yearStr := c.Query("year")
		monthStr := c.Query("month")

		if yearStr == "" || monthStr == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "year and month are required",
			})
		}

		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid year",
			})
		}

		month, err := strconv.Atoi(monthStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid month",
			})
		}

		repo := repository.NewMenuScheduleRepository(db)
		data, err := repo.GetByMonth(year, month)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(data)
	}
}


// ================================
// ADMIN - delete single menu in date
// ================================
func DeleteMenuSchedule(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req struct {
			MenuID int    `json:"menu_id"`
			Date   string `json:"date"`
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		repo := repository.NewMenuScheduleRepository(db)

		if err := repo.DeleteOne(req.MenuID, req.Date); err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"message": "menu removed from date",
		})
	}
}


// ================================
// ADMIN - clear all menus in date
// ================================
func ClearMenuSchedulesByDate(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		date := c.Query("date")
		if date == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "date required",
			})
		}

		repo := repository.NewMenuScheduleRepository(db)

		if err := repo.ClearByDate(date); err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"message": "all menus cleared for date",
		})
	}
}
