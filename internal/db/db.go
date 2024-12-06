package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func New() (*sql.DB, error) {
	driver := os.Getenv("DB_DRIVER")
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// Validate the connection
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
