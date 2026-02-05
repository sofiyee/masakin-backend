package config

import (
	"os"
)

func GetPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	return ":" + port
}
