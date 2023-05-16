package database

func GetAllEntities() []interface{} {
	return []interface{}{
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
	}
}

func AutoMigrate(db GormDB) error {
	err := db.AutoMigrate(
		GetAllEntities()...,
	)
	if err != nil {
		return err
	}

	//https://dba.stackexchange.com/a/164081
	tx := db.Exec(`
CREATE OR REPLACE FUNCTION f_concat_ws(text, VARIADIC text[])
	RETURNS text
	LANGUAGE sql IMMUTABLE PARALLEL SAFE AS
	'SELECT array_to_string($2, $1)';

CREATE INDEX IF NOT EXISTS business_fulltext_idx ON businesses 
	USING GIN (
		to_tsvector('simple', f_concat_ws(' ', name, description, address))
	)`)
	if err := tx.GetError(); err != nil {
		return err
	} else {
		return nil
	}
}
