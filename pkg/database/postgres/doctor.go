package postgres

import (
	"fmt"
)

type Doctor struct {
	Id             int    `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Specialization string `json:"specialization"`
	Room           int    `json:"room"`
	Email          string `json:"email"`
}

func (db *DB) GetAllDoctors() ([]Doctor, error) {
	stmt := `SELECT * FROM doctors`

	rows, err := db.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	doctors := make([]Doctor, 0, 0)

	for rows.Next() {
		doctor := Doctor{}

		err := rows.Scan(&doctor.Id, &doctor.FirstName, &doctor.LastName, &doctor.Specialization, &doctor.Room, &doctor.Email)

		if err != nil {
			return nil, err
		}

		doctors = append(doctors, doctor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil

}

func (db *DB) GetDoctorId(id int) (Doctor, error) {

	stmt := `SELECT * FROM doctors WHERE id=$1`

	row := db.DB.QueryRow(stmt, id)

	d := Doctor{}

	err := row.Scan(&d.Id, &d.FirstName, &d.LastName, &d.Specialization, &d.Room, &d.Email)

	if err != nil {
		return Doctor{}, err
	}
	return d, nil
}

func (db *DB) InsertDoctor(d Doctor) error {
	stmt := `INSERT INTO doctors (firstname, lastname, specialization, email, room) VALUES ($1, $2, $3, $4, $5)`

	if _, err := db.DB.Exec(stmt, d.FirstName, d.LastName, d.Specialization, d.Email, d.Room); err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateDoctorAll(id int, d Doctor) error {

	if _, err := db.GetDoctorId(id); err != nil {
		return err
	}

	stmt := "UPDATE doctors SET  firstname = $1," +
		"lastname = $2, specialization= $3, email = $4, room = $5 WHERE id=$6"

	if _, err := db.DB.Exec(stmt, d.FirstName, d.LastName, d.Specialization, d.Email, d.Room, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateDoctor(id int, req string) error {

	if _, err := db.GetDoctorId(id); err != nil {
		return err
	}

	stmt := fmt.Sprintf("UPDATE doctors SET %s WHERE id=$1", req)

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteDoctor(id int) error {

	if _, err := db.GetDoctorId(id); err != nil {
		return err
	}

	stmt := `DELETE FROM doctors WHERE id=$1`

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) CheckEmailDoctor(email string) error {
	stmt := `SELECT * FROM doctors WHERE email=$1`

	row := db.DB.QueryRow(stmt, email)

	d := Doctor{}

	err := row.Scan(&d.Id, &d.FirstName, &d.LastName, &d.Specialization, &d.Room, &d.Email)
	return err
}
