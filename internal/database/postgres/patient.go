package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Patient struct {
	Id           *int    `json:"id"`
	FirstName    *string `json:"firstName"`
	LastName     *string `json:"lastName"`
	Email        *string `json:"email"` // may be null in DB
	Address      *string `json:"address"`
	Phone_Number *string `json:"phone_number"`
	Login        *string `json:"login"`
	Password     *string `json:"password"`
	RoleId       *int    `json:"role_id"`
	RefreshToken *string `json:"refreshToken"`
}

func (p *Patient) DecodeFromJSON(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(p); err != nil {
		return err
	}
	return nil
}

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
			&patient.Phone_Number, &patient.Login, &patient.Password, &patient.RoleId, &patient.RefreshToken)

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
	stmt := `SELECT * FROM patients WHERE id=$1 LIMIT 1;`

	row := db.DB.QueryRow(stmt, id)

	p := Patient{}

	err := row.Scan(&p.Id, &p.FirstName, &p.LastName, &p.Email, &p.Address,
		&p.Phone_Number, &p.Login, &p.Password, &p.RoleId, &p.RefreshToken)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (db *DB) InsertPatient(p Patient) error {
	stmt := `INSERT INTO patients (firstname, lastname, email, address, phone_number,
                      login, password, role_id, refresh_token) VALUES ($1, $2, $3, $4,$5,$6,$7,$8,$9);`

	if _, err := db.DB.Exec(stmt, p.FirstName, p.LastName, p.Email, p.Address,
		p.Phone_Number, p.Login, p.Password, p.RoleId, p.RefreshToken); err != nil {
		return err
	}

	return nil
}

//func (db *DB) UpdatePatientAll(id int, p Patient) error {
//	if _, err := db.GetPatientId(id); err != nil {
//		return err
//	}
//
//	stmt := "UPDATE patients SET  firstname = $1," +
//		"lastname = $2, email= $3, address = $4, phone_number = $5 WHERE id=$6;"
//
//	if _, err := db.DB.Exec(stmt, p.FirstName, p.LastName, p.Email, p.Address, p.Phone_Number, id); err != nil {
//		return err
//	}
//
//	return nil
//}

func (db *DB) UpdatePatient(id int, req string) error {
	if _, err := db.GetPatientId(id); err != nil {
		return err
	}

	stmt := fmt.Sprintf("UPDATE patients  SET %s WHERE id=$1;", req)

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeletePatient(id int) error {

	if _, err := db.GetPatientId(id); err != nil {
		return err
	}

	stmt := `DELETE FROM patients WHERE id=$1;`

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) CheckEmailPatient(email string) error {
	stmt := fmt.Sprintf(`SELECT email FROM patients WHERE email=$1 LIMIT 1;`)

	row := db.DB.QueryRow(stmt, email)

	var patientEmail string
	err := row.Scan(&patientEmail)

	return err
}

func (db *DB) CheckPhoneNumberPatient(phoneNumber string) error {
	stmt := fmt.Sprintf(`SELECT * FROM patients WHERE phone_number=$1 LIMIT 1;`)

	row := db.DB.QueryRow(stmt, phoneNumber)

	var patientPhone string
	err := row.Scan(&patientPhone)

	return err
}

func (db *DB) CheckLoginAndPasswordPatient(login, password string) (int, error) {
	stmt := `SELECT id FROM patients WHERE login=$1 AND password=$2 LIMIT 1;`

	row := db.DB.QueryRow(stmt, login, password)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *DB) UpdatePatientRefreshToken(id int, refreshToken string) error {
	stmt := "UPDATE patients  SET refresh_token=$1 WHERE id=$2;"

	if _, err := db.DB.Exec(stmt, refreshToken, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) CheckLoginPatient(login string) error {
	stmt := `SELECT id FROM patients WHERE login=$1 LIMIT 1;`

	row := db.DB.QueryRow(stmt, login)

	var id int
	err := row.Scan(&id)

	return err
}

func (db *DB) CheckPatientRefreshToken(id int) (string, error) {
	stmt := `SELECT refresh_token FROM patients WHERE id=$1 LIMIT 1;`

	row := db.DB.QueryRow(stmt, id)

	var refreshToken string
	err := row.Scan(&refreshToken)

	return refreshToken, err
}
