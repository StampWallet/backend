package database

func AutoMigrate(db GormDB) error {
	err := db.AutoMigrate(
		&User{},
		&LocalCard{},
		&Token{},
		&Business{},
		&FileMetadata{},
		&ItemDefinition{},
		&MenuItem{},
		&OwnedItem{},
		&VirtualCard{},
		&Transaction{},
		&TransactionDetail{},
	)
	if err != nil {
		return err
	}

	tx := db.Exec("CREATE INDEX IF NOT EXISTS business_fulltext_idx ON businesses USING GIN (to_tsvector('simple', name || ' ' || description || ' ' || address))")
	if err := tx.GetError(); err != nil {
		return err
	} else {
		return nil
	}
}
