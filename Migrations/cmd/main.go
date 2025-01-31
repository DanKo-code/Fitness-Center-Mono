package main

import (
	"Migrations/constants"
	logrusCustom "Migrations/pkg/logger"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logrusCustom.InitLogger()

	err := godotenv.Load()
	if err != nil {
		logrusCustom.LogWithLocation(logrus.FatalLevel, fmt.Sprintf("Error loading .env file: %s", err))
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, "Successfully loaded environment variables")

	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_DRIVER"),

		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),

		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),

		os.Getenv("DB_NAME"),

		os.Getenv("DB_SLLMODE"),
	)

	db, err := sql.Open(os.Getenv("DB_DRIVER"), dsn)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.FatalLevel, fmt.Sprintf("Failed to create DB object: %v", err))
	}

	isFirstEnter := checkFirstEnter(db)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logrusCustom.LogWithLocation(logrus.FatalLevel, fmt.Sprintf("Failed to create migration driver: %v", err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		os.Getenv("DB_NAME"),
		driver)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.FatalLevel, fmt.Sprintf("Failed to initialize migration: %v", err))
	}

	if !isFirstEnter {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, "Migrations table exists. Skipping migrations.")
		return
	}

	if err := m.Up(); err != nil {
		logrusCustom.LogWithLocation(logrus.FatalLevel, fmt.Sprintf("Failed to run migrations: %v", err))
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, "Successfully run migrations")

	initData(db)
}

func checkFirstEnter(db *sql.DB) bool {
	var version int

	err := db.QueryRow("SELECT version FROM schema_migrations").Scan(&version)
	if err != nil {

		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if string(pqError.Code) == constants.RelationDoesNotExistConst {
				return true
			}
		}

		return false
	}

	if version != 0 {
		return false
	}

	// If down under first migration to null "version" field on "schema_migrations" table
	return false
}

func initData(db *sql.DB) {

	file, err := os.ReadFile(os.Getenv("DB_INIT_DATA_SCRIPT"))
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, fmt.Sprintf("Failed to read DB_INIT_DATA_SCRIPT file: %v", err))
		return
	}

	_, err = db.Exec(string(file))
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, fmt.Sprintf("Error executing DB_INIT_DATA_SCRIPT: %v", err))
		return
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, "Successfully insert init data")
}
