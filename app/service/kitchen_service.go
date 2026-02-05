package service

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"masakin-backend/app/repository"
	"masakin-backend/app/model"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func CreateGudangUser(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req struct {
			Name     string `json:"name"`
			Phone    string `json:"phone"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		if req.Name == "" || req.Phone == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "name, phone, password are required",
			})
		}

		hashed, err := bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		userRepo := repository.NewUserRepository(db)

		_, err = userRepo.Create(&model.User{
			Name:     req.Name,
			Phone:    req.Phone,
			Address:  "Dapur",
			Role:     "gudang",
			Password: string(hashed),
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "gudang user created",
		})
	}
}


func GetDailyKitchenSummary(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// 🔍 1. CEK ROLE & USER (PALING PENTING)
		log.Println(
			"KITCHEN ACCESS:",
			"user_id =", c.Locals("user_id"),
			"role =", c.Locals("role"),
		)

		// 🔍 2. CEK QUERY PARAM
		date := c.Query("date") // YYYY-MM-DD
		log.Println("KITCHEN DATE:", date)

		if date == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "date is required (YYYY-MM-DD)",
			})
		}

		repo := repository.NewKitchenRepository(db)

		// 🔍 3. CEK ERROR QUERY DB
		data, err := repo.GetDailySummary(date)
		if err != nil {
			log.Println("KITCHEN QUERY ERROR:", err)
			return fiber.ErrInternalServerError
		}

		log.Println("KITCHEN RESULT COUNT:", len(data))

		return c.JSON(fiber.Map{
			"date":  date,
			"menus": data,
		})
	}
}
