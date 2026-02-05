package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"masakin-backend/config"
	"masakin-backend/database"
	"masakin-backend/routes"
	
)

func main() {
	config.LoadEnv()
	db := database.ConnectPostgres()

	app := fiber.New()

	// 🔥 CORS FIX
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	app.Static("/uploads", "./uploads")

	routes.Register(app, db)

	log.Println("🚀 Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

