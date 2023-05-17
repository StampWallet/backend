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

func RecreateDatabase(db GormDB, databaseName string) error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 100 * time.Second,
		},
	)

	//https://dba.stackexchange.com/a/154075
	tx := db.Session(&gorm.Session{Logger: newLogger}).Exec(`
do
$$
declare
  l_stmt text;
begin
  select 'truncate ' || string_agg(format('%I.%I', schemaname, tablename), ',')
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
		return err
	}
	if err := AutoMigrate(db); err != nil {
		return err
	}
	return nil
}

func GetTestDatabase() *GormDBImpl {
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
	db.Exec("SELECT 'Migration finished'")
	return &GormDBImpl{Db: db}
}
