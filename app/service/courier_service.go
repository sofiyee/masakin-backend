package service

import (
	"database/sql"
	"log"
	"masakin-backend/app/model"
	"masakin-backend/app/repository"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"path/filepath"
	"fmt"
	"time"
	"os"
	"strings"

)

// CREATE COURIER

func CreateCourier(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req struct {
			Name     string `json:"name"`
			Phone    string `json:"phone"`
			Address  string `json:"address"`
			Password string `json:"password"`
			Region   string `json:"region"`
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		userRepo := repository.NewUserRepository(db)
		courierRepo := repository.NewCourierRepository(db)

		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		user := &model.User{
			Name:     req.Name,
			Phone:    req.Phone,
			Address:  req.Address,
			Password: string(hash),
			Role:     "kurir",
		}

		userID, err := userRepo.Create(user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if err := courierRepo.Create(userID, req.Region); err != nil {
			return fiber.ErrInternalServerError
		}

		log.Printf("🚚 Courier created: %s (%s)", req.Name, req.Region)

		return c.JSON(fiber.Map{"message": "courier created"})
	}
}

// LIST COURIERS

func GetCouriers(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		courierRepo := repository.NewCourierRepository(db)

		data, err := courierRepo.GetAll()
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(data)
	}
}

// UPDATE COURIER REGION

func UpdateCourierRegion(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req struct {
			UserID int    `json:"user_id"`
			Region string `json:"region"`
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		repo := repository.NewCourierRepository(db)
		if err := repo.UpdateRegion(req.UserID, req.Region); err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{"message": "region updated"})
	}
}

// SET COURIER ACTIVE STATUS
func SetCourierActive(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req struct {
			UserID int  `json:"user_id"`
			Active bool `json:"active"`
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		repo := repository.NewCourierRepository(db)
		if err := repo.SetActive(req.UserID, req.Active); err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{"message": "courier status updated"})
	}
}
// ================================
// GET ORDERS FOR COURIER
// ================================
func GetOrdersForCourier(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		kurirID := c.Locals("user_id").(int)

		orderRepo := repository.NewOrderRepository(db)
		orders, err := orderRepo.GetForCourier(kurirID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(orders)
	}
}

// ================================
// DELIVER ORDER (KURIR + FOTO)
// ================================
func DeliverOrder(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// =====================
		// 0. AUTH CONTEXT
		// =====================
		userIDRaw := c.Locals("user_id")
		if userIDRaw == nil {
			log.Println("AUTH ERROR: user_id not found in context")
			return fiber.ErrUnauthorized
		}
		userID := userIDRaw.(int) // ini USERS.ID

		orderID, err := strconv.Atoi(c.Params("id"))
		if err != nil || orderID == 0 {
			log.Println("PARAM ERROR: invalid order id:", c.Params("id"), "err:", err)
			return fiber.ErrBadRequest
		}

		// =====================
		// 1. AMBIL COURIER.ID (INI KUNCI FIX FK ❗)
		// =====================
		var courierID int
		err = db.QueryRow(`
			SELECT id
			FROM couriers
			WHERE user_id = $1
			AND active = true
		`, userID).Scan(&courierID)

		if err != nil {
			log.Println("COURIER LOOKUP ERROR:", err, "user_id:", userID)
			return fiber.ErrForbidden
		}

		// =====================
		// 2. FILE FOTO
		// =====================
		file, err := c.FormFile("proof_image")
		if err != nil {
			log.Println("FORM FILE ERROR: proof_image missing:", err)
			return c.Status(400).JSON(fiber.Map{
				"error": "proof_image is required",
			})
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowed := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		if !allowed[ext] {
			log.Println("FILE TYPE ERROR:", file.Filename)
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid image type",
			})
		}

		filename := fmt.Sprintf(
			"delivery_%d_%d%s",
			orderID,
			time.Now().Unix(),
			ext,
		)

		saveDir := "./uploads/deliveries"
		savePath := saveDir + "/" + filename

		// pastikan folder ada
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			log.Println("MKDIR ERROR:", err)
			return fiber.ErrInternalServerError
		}

		if err := c.SaveFile(file, savePath); err != nil {
			log.Println("SAVE FILE ERROR:", err, "path:", savePath)
			return fiber.ErrInternalServerError
		}

		// =====================
		// 3. VALIDASI ORDER (REGION + STATUS)
		// =====================
		var valid bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM orders o
				JOIN customers c ON c.id = o.customer_id
				JOIN couriers cr ON cr.region = c.region
				WHERE o.id = $1
				AND cr.id = $2
				AND cr.active = true
				AND o.status = 'paid'
			)
		`, orderID, courierID).Scan(&valid)

		if err != nil {
			log.Println("VALIDATION QUERY ERROR:", err)
			return fiber.ErrInternalServerError
		}

		if !valid {
			log.Println(
				"FORBIDDEN: order not valid for courier",
				"orderID:", orderID,
				"courierID:", courierID,
			)
			return fiber.ErrForbidden
		}

		// =====================
		// 4. TRANSACTION
		// =====================
		tx, err := db.Begin()
		if err != nil {
			log.Println("TX BEGIN ERROR:", err)
			return fiber.ErrInternalServerError
		}
		defer tx.Rollback()

		deliveryRepo := repository.NewDeliveryRepository(tx)
		orderRepo := repository.NewOrderRepository(tx)

		// insert delivery
		if err := deliveryRepo.Create(orderID, courierID, filename); err != nil {
			log.Println("INSERT DELIVERY ERROR:", err)
			return fiber.ErrInternalServerError
		}

		// update order status
		if err := orderRepo.Deliver(orderID); err != nil {
			log.Println("UPDATE ORDER ERROR:", err)
			return fiber.ErrInternalServerError
		}

		if err := tx.Commit(); err != nil {
			log.Println("TX COMMIT ERROR:", err)
			return fiber.ErrInternalServerError
		}

		log.Println(
			"DELIVERY SUCCESS",
			"orderID:", orderID,
			"courierID:", courierID,
			"file:", filename,
		)

		return c.JSON(fiber.Map{
			"message": "order delivered with proof",
		})
	}
}


