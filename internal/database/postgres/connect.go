package postgres

import (
	"database/sql"
)

type DB struct {
	DB *sql.DB
}

func ConnectDB(infoStr string) (*sql.DB, error) {

	db, err := sql.Open("postgres", infoStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
