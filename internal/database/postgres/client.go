package postgres

import (
	"fmt"
	"strings"
)

type Client struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Address   string `json:"address"`
}

var CanChangeFieldClient int = 4

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

	stmt := `SELECT * FROM clients WHERE id=$1`

	rows := db.DB.QueryRow(stmt, id)

	c := &Client{}

	err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Email, &c.Address)

	if err != nil {
		return nil, err
	}
	return c, nil
}

func (db *DB) InsertClient(c *Client) error {
	stmt := `INSERT INTO clients (firstname, lastname, email, address) VALUES ($1, $2, $3, $4)`

	if _, err := db.DB.Exec(stmt, c.FirstName, c.LastName, c.Email, c.Address); err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateClientAll(id int, c *Client) error {

	if _, err := db.GetClient(id); err != nil {
		return err
	}

	stmt := "UPDATE clients SET  firstname = $1," +
		"lastname = $2, email= $3, address = $4 WHERE id=$5"

	if _, err := db.DB.Exec(stmt, c.FirstName, c.LastName, c.Email, c.Address, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateClient(id int, m map[string]string) error {

	if _, err := db.GetClient(id); err != nil {
		return err
	}

	var b strings.Builder

	for key, val := range m {
		b.WriteString(fmt.Sprintf("%s='%s',", key, val))
	}

	str := b.String()

	str = str[:len(str)-1]

	stmt := fmt.Sprintf("UPDATE clients SET %s WHERE id=$1", str)

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteClient(id int) error {

	if _, err := db.GetClient(id); err != nil {
		return err
	}

	stmt := fmt.Sprintf(`DELETE FROM clients WHERE id=$1`)

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}
