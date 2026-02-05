package service

import (
	"masakin-backend/app/model"
	"masakin-backend/app/repository"
	"github.com/gofiber/fiber/v2"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"masakin-backend/utils"
	"log"
	"strings"

)

type AuthService struct {
	UserRepo     *repository.UserRepository
	CustomerRepo *repository.CustomerRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
	customerRepo *repository.CustomerRepository,
) *AuthService {
	return &AuthService{
		UserRepo:     userRepo,
		CustomerRepo: customerRepo,
	}
}

// LOGIN (SEMUA ROLE)

func Login(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req struct {
			Name          string `json:"name"`
			Password      string `json:"password"`
			CaptchaID     string `json:"captcha_id"`
			CaptchaAnswer int    `json:"captcha_answer"`
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		// Captcha validation
		if !VerifyMathCaptcha(req.CaptchaID, req.CaptchaAnswer) {
			return c.Status(400).JSON(fiber.Map{
				"error": "captcha invalid",
			})
		}

		userRepo := repository.NewUserRepository(db)

		user, err := userRepo.FindByName(req.Name)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "user not found",
			})
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(req.Password),
		); err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid credentials",
			})
		}

		token, err := utils.GenerateToken(user.ID, user.Name, user.Role)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to generate token",
			})
		}

		log.Printf("✅ Login success: %s (%s)", user.Name, user.Role)

		return c.JSON(fiber.Map{
			"token": token,
			"user": fiber.Map{
				"name": user.Name,
				"role": user.Role,
			},
		})
	}
}

// Logout
func Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {

		name, _ := c.Locals("name").(string)
		role, _ := c.Locals("role").(string)

		log.Printf("🚪 Logout success: %s (%s)", name, role)

		return c.JSON(fiber.Map{
			"message": "logout successful",
		})
	}
}



// REGISTER CUSTOMER
func RegisterCustomer(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// =====================
		// 1. PARSE REQUEST
		// =====================
		var req struct {
			Name        string `json:"name"`
			Phone       string `json:"phone"`
			Region      string `json:"region"`
			FullAddress string `json:"full_address"`
			Password    string `json:"password"`
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		// =====================
		// 2. VALIDASI WAJIB
		// =====================
		if req.Name == "" ||
			req.Phone == "" ||
			req.Region == "" ||
			req.FullAddress == "" ||
			req.Password == "" {

			return c.Status(400).JSON(fiber.Map{
				"error": "name, phone, region, full_address, and password are required",
			})
		}

		userRepo := repository.NewUserRepository(db)
		customerRepo := repository.NewCustomerRepository(db)

		// =====================
		// 3. HASH PASSWORD
		// =====================
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		// =====================
		// 4. CREATE USER (AUTH)
		// =====================
		user := &model.User{
			Name:     req.Name,
			Phone:    req.Phone,
			Password: string(hashedPassword),
			Role:     "customer",
		}

		userID, err := userRepo.Create(user)
		if err != nil {
			if strings.Contains(err.Error(), "users_phone_unique") {
				return c.Status(409).JSON(fiber.Map{
					"error": "phone number already registered",
				})
			}

			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// =====================
		// 5. CREATE CUSTOMER PROFILE
		// =====================
		err = customerRepo.Create(&model.Customer{
			UserID:      userID,
			Name:        req.Name,
			Region:      req.Region,
			FullAddress: req.FullAddress,
		})
		if err != nil {
			return fiber.ErrInternalServerError
		}

		// =====================
		// 6. RESPONSE
		// =====================
		return c.JSON(fiber.Map{
			"message": "register success",
		})
	}
}






