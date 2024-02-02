package postgres

import (
	"fmt"
	"strings"
)

type Doctor struct {
	Id             int    `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Specialization string `json:"specialization"`
	Room           int    `json:"room"`
}

var CanChangeFieldDoctor int = 4

func (db *DB) GetAllDoctors() ([]*Doctor, error) {
	stmt := `SELECT * FROM doctors`

	rows, err := db.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	doctors := make([]*Doctor, 0, 0)

	for rows.Next() {
		doctor := &Doctor{}

		rows.Scan(&doctor.Id, &doctor.FirstName, &doctor.LastName, &doctor.Specialization, &doctor.Room)

		doctors = append(doctors, doctor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil

}

func (db *DB) GetDoctor(id int) (*Doctor, error) {

	stmt := `SELECT * FROM doctors WHERE id=$1`

	rows := db.DB.QueryRow(stmt, id)

	d := &Doctor{}

	err := rows.Scan(&d.Id, &d.FirstName, &d.LastName, &d.Specialization, &d.Room)

	if err != nil {
		return nil, err
	}
	return d, nil
}

func (db *DB) InsertDoctor(d *Doctor) error {
	stmt := `INSERT INTO doctors (firstname, lastname, specialization, room) VALUES ($1, $2, $3, $4)`

	if _, err := db.DB.Exec(stmt, d.FirstName, d.LastName, d.Specialization, d.Room); err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateDoctorAll(id int, d *Doctor) error {

	if _, err := db.GetDoctor(id); err != nil {
		return err
	}

	stmt := "UPDATE doctors SET  firstname = $1," +
		"lastname = $2, specialization= $3, room = $4 WHERE id=$5"

	if _, err := db.DB.Exec(stmt, d.FirstName, d.LastName, d.Specialization, d.Room, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateDoctor(id int, m map[string]string) error {

	if _, err := db.GetDoctor(id); err != nil {
		return err
	}

	var b strings.Builder

	for key, val := range m {

		if key == "room" {
			b.WriteString(fmt.Sprintf("%s=%s,", key, val))
			continue
		}
		b.WriteString(fmt.Sprintf("%s='%s',", key, val))
	}

	str := b.String()

	str = str[:len(str)-1]

	stmt := fmt.Sprintf("UPDATE doctors SET %s WHERE id=$1", str)

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteDoctor(id int) error {

	if _, err := db.GetDoctor(id); err != nil {
		return err
	}

	stmt := fmt.Sprintf(`DELETE FROM doctors WHERE id=$1`)

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}
