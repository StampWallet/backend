package testutils

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	. "github.com/StampWallet/backend/internal/database"
)

// Wipes the database
func RecreateDatabase(db GormDB, databaseName string) error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 100 * time.Second,
		},
	)

	for i := 0; i <= 5; i++ {
		// Truncates all tables
		//https://dba.stackexchange.com/a/154075
		tx := db.Session(&gorm.Session{Logger: newLogger}).Exec(`
	do
	$$
	declare
	  l_stmt text;
	begin
	  select 'truncate ' || string_agg(format('%I.%I', schemaname, tablename), ',') || ' cascade'
		into l_stmt
	  from pg_tables
	  where schemaname in ('public') and tablename != 'spatial_ref_sys';

	  if l_stmt is not null then 
		execute l_stmt;
	  end if;
	end;
	$$
		`)
		if err := tx.GetError(); err != nil {
			if i == 5 {
				return err
			} else {
				time.Sleep(time.Second * 10)
			}
		}
	}

	// Applies auto migration
	if err := AutoMigrate(db); err != nil {
		return err
	}
	return nil
}

// Creates a new database connection from environment variables. Wipes the database.
// TEST_DATABASE_URL - database URL, ex. 'postgres://postgres@localhost/stampwallet'
// TEST_DATABASE_NAME - database name, ex. 'stampwallet'
func GetTestDatabase() *GormDBImpl {
	// Get environment variables
	url := os.Getenv("TEST_DATABASE_URL")
	dbname := os.Getenv("TEST_DATABASE_NAME")
	if url == "" {
		panic("no database url in TEST_DATABASE_URL env var")
	}
	if dbname == "" {
		panic("no database name in TEST_DATABASE_NAME env var")
	}

	// Create db connection
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: url,
	}))
	if err != nil {
		panic(fmt.Errorf("Failed to open database connection %w", err))
	}

	// Recreate database
	if err := RecreateDatabase(&GormDBImpl{Db: db}, os.Getenv("TEST_DATABASE_NAME")); err != nil {
		panic(fmt.Errorf("failed to drop and recreate database %w", err))
	}
	db.Exec("SELECT 'Migration finished'")
	return &GormDBImpl{Db: db}
}
