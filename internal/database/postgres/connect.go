package postgres

import (
	"database/sql"
)

type DB struct {
	DB *sql.DB
}

var DataBase = DB{}

func InitDB(infoStr string) (*sql.DB, error) {

	db, err := sql.Open("postgres", infoStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
