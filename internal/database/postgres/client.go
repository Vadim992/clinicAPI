package postgres

import (
	"fmt"
)

type Client struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Address   string `json:"address"`
}

func (db *DB) GetAllClients() ([]*Client, error) {
	stmt := `SELECT * FROM clients`

	rows, err := db.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := make([]*Client, 0, 0)

	for rows.Next() {
		client := &Client{}

		rows.Scan(&client.Id, &client.FirstName, &client.LastName, &client.Email, &client.Address)

		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil

}

func (db *DB) GetClient(id int) (*Client, error) {

	stmt := fmt.Sprintf(`SELECT * FROM clients WHERE id=$1`)

	rows := db.DB.QueryRow(stmt, id)

	c := &Client{}

	err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Email, &c.Address)

	if err != nil {
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
