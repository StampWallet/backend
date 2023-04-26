package database

func AutoMigrate(db GormDB) error {
	return db.AutoMigrate(
		&User{},
		&LocalCard{},
		&Token{},
		&FileMetadata{},
		&Business{},
		&ItemDefinition{},
		&MenuItem{},
		&OwnedItem{},
		&VirtualCard{},
		&Transaction{},
		&TransactionDetails{},
	)
}
