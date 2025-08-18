package nginx_parser

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "gorm.io/driver/postgres"
)

func NewDb() (db *sql.DB, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")

	fmt.Printf("Connecting to database: %s\n", dsn)
	db, err = sql.Open("pgx", string(dsn))

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, nil
}
