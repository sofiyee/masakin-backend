package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"
)


// ============================
// POSTGRESQL CONNECTION
// ============================
var DB *sql.DB

func ConnectPostgres() *sql.DB {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("❌ DB_DSN tidak ditemukan di .env")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Gagal konek ke PostgreSQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ PostgreSQL tidak bisa di-ping: %v", err)
	}

	DB = db
	fmt.Println("✅ PostgreSQL berhasil terkoneksi!")
	return db
}