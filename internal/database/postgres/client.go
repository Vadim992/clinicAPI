package postgres

import (
	"fmt"
)

type Client struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Address   string
}

func (db *DB) GetClient(id int) (*Client, error) {

	stmt := fmt.Sprintf(`SELECT * FROM clients WHERE id=$1`)

	rows := db.DB.QueryRow(stmt, id)

	c := &Client{}

	if err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Email, &c.Address); err != nil {
		return nil, err
	}
	return c, nil
}

func (db *DB) InsertData(firstName, lastName, email, address string) error {
	stmt := fmt.Sprintf("INSERT INTO clients (firstname, lastname, email, address) VALUES ($1, $2, $3, $4)")

	if _, err := db.DB.Exec(stmt, firstName, lastName, email, address); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteData(id int) error {
	stmt := fmt.Sprintf("DELETE FROM clients WHERE id=$1")

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) Update(id int) error {
	stmt := fmt.Sprintf("UPDATE clients SET id=5 WHERE id=$1")

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}
