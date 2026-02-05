package main

import (
	"log"
	"os"

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

	// 🔥 CORS (dev friendly)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // sementara, biar FE temenmu bisa akses
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	app.Static("/uploads", "./uploads")

	routes.Register(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server running on port %s\n", port)
	log.Fatal(app.Listen(":" + port))
}
