package postgres

import (
	"fmt"
)

func (db *DB) GetPatients(offset, limit int, filter string) ([]Patient, error) {
	var order string

	if filter != "" {
		order = fmt.Sprintf("ORDER BY %s", filter)
	}
	stmt := fmt.Sprintf(`SELECT * FROM patients %s
    LIMIT $1 OFFSET $2;`, order)

	rows, err := db.DB.Query(stmt, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patients := make([]Patient, 0, 0)

	for rows.Next() {
		patient := Patient{}

		err := rows.Scan(&patient.Id, &patient.FirstName, &patient.LastName, &patient.Email, &patient.Address,
			&patient.Phone_Number)

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

	err := row.Scan(&p.Id, &p.FirstName, &p.LastName, &p.Email, &p.Address,
		&p.Phone_Number)

	if err != nil {
		return p, err
	}
	return p, nil
}

func (db *DB) InsertPatient(p Patient) error {

	stmt := `INSERT INTO patients (firstname, lastname, email, address, phone_number) VALUES ($1, $2, $3, $4,$5)`

	if _, err := db.DB.Exec(stmt, p.FirstName, p.LastName, p.Email, p.Address, &p.Phone_Number); err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdatePatientAll(id int, p Patient) error {

	if _, err := db.GetPatientId(id); err != nil {
		return err
	}

	stmt := "UPDATE patients SET  firstname = $1," +
		"lastname = $2, email= $3, address = $4, phone_number = $5 WHERE id=$6"

	if _, err := db.DB.Exec(stmt, p.FirstName, p.LastName, p.Email, p.Address, p.Phone_Number, id); err != nil {
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

func (db *DB) CheckPhoneNumberPatient(phoneNumber string) error {
	stmt := fmt.Sprintf(`SELECT * FROM patients WHERE phone_number=$1`)

	row := db.DB.QueryRow(stmt, phoneNumber)

	p := Patient{}
	err := row.Scan(&p.Id, &p.FirstName, &p.LastName, &p.Email, &p.Address, &p.Phone_Number)

	return err
}
