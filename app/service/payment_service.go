package service
import (
	"database/sql"
	"masakin-backend/app/repository"
	"github.com/gofiber/fiber/v2"
	"fmt"
	"path/filepath"
	"strconv"
	"time"
	"os"
)

func CreatePayment(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// =====================
		// 1. ADMIN ID
		// =====================
		adminIDRaw := c.Locals("user_id")
		if adminIDRaw == nil {
			return fiber.ErrUnauthorized
		}
		adminID := adminIDRaw.(int)

		// =====================
		// 2. ORDER ID (FORM)
		// =====================
		orderIDStr := c.FormValue("order_id")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil || orderID == 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "order_id is required",
			})
		}

		// =====================
		// 3. FILE UPLOAD
		// =====================
		file, err := c.FormFile("proof_image")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "proof_image is required",
			})
		}

		ext := filepath.Ext(file.Filename)
		allowed := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		if !allowed[ext] {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid file type",
			})
		}

		filename := fmt.Sprintf(
			"payment_%d_%d%s",
			orderID,
			time.Now().Unix(),
			ext,
		)

		savePath := "./uploads/payments/" + filename
		if err := c.SaveFile(file, savePath); err != nil {
			return fiber.ErrInternalServerError
		}

		// =====================
		// 4. TRANSACTION
		// =====================
		tx, err := db.Begin()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		defer tx.Rollback()

		paymentRepo := repository.NewPaymentRepository(tx)
		orderRepo   := repository.NewOrderRepository(tx)

		if err := paymentRepo.Create(orderID, filename, adminID); err != nil {
			return fiber.ErrInternalServerError
		}

		if err := orderRepo.UpdateStatus(orderID, "paid"); err != nil {
			return fiber.ErrInternalServerError
		}

		if err := tx.Commit(); err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"message": "payment recorded & order marked as paid",
		})
	}
}

func UpdatePaymentImage(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// =====================
		// 1. ADMIN AUTH
		// =====================
		adminIDRaw := c.Locals("user_id")
		if adminIDRaw == nil {
			return fiber.ErrUnauthorized
		}

		// =====================
		// 2. ORDER ID
		// =====================
		orderIDStr := c.FormValue("order_id")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil || orderID == 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "order_id is required",
			})
		}

		// =====================
		// 3. FILE
		// =====================
		file, err := c.FormFile("proof_image")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "proof_image is required",
			})
		}

		ext := filepath.Ext(file.Filename)
		allowed := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		if !allowed[ext] {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid file type",
			})
		}

		// =====================
		// 4. TRANSACTION
		// =====================
		tx, err := db.Begin()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		defer tx.Rollback()

		paymentRepo := repository.NewPaymentRepository(tx)

		// ambil file lama
		oldFile, err := paymentRepo.GetProofImageByOrderID(orderID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "payment not found",
			})
		}

		// simpan file baru
		newFilename := fmt.Sprintf(
			"payment_%d_%d%s",
			orderID,
			time.Now().Unix(),
			ext,
		)
		savePath := "./uploads/payments/" + newFilename

		if err := c.SaveFile(file, savePath); err != nil {
			return fiber.ErrInternalServerError
		}

		// update db
		if err := paymentRepo.UpdateProofImage(orderID, newFilename); err != nil {
			return fiber.ErrInternalServerError
		}

		// commit dulu sebelum hapus file lama
		if err := tx.Commit(); err != nil {
			return fiber.ErrInternalServerError
		}

		// =====================
		// 5. HAPUS FILE LAMA
		// =====================
		if oldFile != "" {
			_ = os.Remove("./uploads/payments/" + oldFile)
		}

		return c.JSON(fiber.Map{
			"message": "payment image updated",
		})
	}
}

func GetAllPaymentsAdmin(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		repo := repository.NewPaymentRepository(db)

		payments, err := repo.GetAllAdmin()
		if err != nil {
			fmt.Println("GET PAYMENTS ERROR:", err)
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(payments)
	}
}

