package postgres

import (
	"database/sql"
	"fmt"
)

type Patient struct {
	Id        int            `json:"id"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Email     sql.NullString `json:"email"`
	Address   string         `json:"address"`
}

func (db *DB) GetPatients() ([]Patient, error) {
	stmt := `SELECT * FROM patients `

	rows, err := db.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patients := make([]Patient, 0, 0)

	for rows.Next() {
		patient := Patient{}

		err := rows.Scan(&patient.Id, &patient.FirstName, &patient.LastName, &patient.Email, &patient.Address)

		if err != nil {
			return nil, err
		}

		patients = append(patients, patient)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil

}

func (db *DB) GetPatientId(id int) (Patient, error) {

	stmt := `SELECT * FROM patients WHERE id=$1`

	row := db.DB.QueryRow(stmt, id)

	p := Patient{}

	err := row.Scan(&p.Id, &p.FirstName, &p.LastName, &p.Email, &p.Address)

	if err != nil {
		return p, err
	}
	return p, nil
}

func (db *DB) InsertPatient(p Patient) error {

	stmt := `INSERT INTO patients (firstname, lastname, email, address) VALUES ($1, $2, $3, $4)`

	if _, err := db.DB.Exec(stmt, p.FirstName, p.LastName, p.Email, p.Address); err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdatePatientAll(id int, p Patient) error {

	if _, err := db.GetPatientId(id); err != nil {
		return err
	}

	stmt := "UPDATE patients SET  firstname = $1," +
		"lastname = $2, email= $3, address = $4 WHERE id=$5"

	if _, err := db.DB.Exec(stmt, p.FirstName, p.LastName, p.Email, p.Address, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdatePatient(id int, req string) error {

	if _, err := db.GetPatientId(id); err != nil {
		return err
	}

	stmt := fmt.Sprintf("UPDATE patients  SET %s WHERE id=$1", req)

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeletePatient(id int) error {

	if _, err := db.GetPatientId(id); err != nil {
		return err
	}

	stmt := `DELETE FROM patients WHERE id=$1`

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) CheckEmailPatient(email string) error {
	stmt := fmt.Sprintf(`SELECT * FROM patients WHERE email=$1`)

	row := db.DB.QueryRow(stmt, email)

	p := Patient{}
	err := row.Scan(&p.Id, &p.FirstName, &p.LastName, &p.Email, &p.Address)

	return err
}
