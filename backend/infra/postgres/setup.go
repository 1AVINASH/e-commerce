package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Inititialize() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	sslmode := os.Getenv("POSTGRES_SSLMODE") // e.g., "disable"

	if host == "" || port == "" || user == "" || dbname == "" {
		log.Fatal("Postgres configuration missing in environment variables")
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening Postgres: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Check connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Error pinging Postgres: %v", err)
	}

	log.Println("✅ Connected to Postgres successfully")
}

func ClosePostgres() {
	if DB != nil {
		DB.Close()
		log.Println("Postgres connection closed")
	}
}
