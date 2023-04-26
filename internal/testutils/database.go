package testutils

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
)

func RecreateDatabase(db GormDB, databaseName string) error {
	tx := db.Raw("DROP SCHEMA ? CASCADE; CREATE SCHEMA ?;", databaseName, databaseName)
	if err := tx.GetError(); err != nil {
		return err
	}
	if err := AutoMigrate(db); err != nil {
		return err
	}
	return nil
}

func GetDatabase() *GormDBImpl {
	url := os.Getenv("TEST_DATABASE_URL")
	dbname := os.Getenv("TEST_DATABASE_NAME")
	if url == "" {
		panic("no database url in TEST_DATABASE_URL env var")
	}
	if dbname == "" {
		panic("no database name in TEST_DATABASE_NAME env var")
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: url,
	}))
	if err != nil {
		panic(fmt.Errorf("Failed to open database connection %w", err))
	}
	if err := RecreateDatabase(&GormDBImpl{Db: db}, os.Getenv("TEST_DATABASE_NAME")); err != nil {
		panic(fmt.Errorf("failed to drop and recreate database %w", err))
	}
	db.Exec("SELECT 1")
	return &GormDBImpl{Db: db}
}
