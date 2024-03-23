package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Doctor struct {
	Id             *int    `json:"id"`
	FirstName      *string `json:"firstName"`
	LastName       *string `json:"lastName"`
	Specialization *string `json:"specialization"`
	Room           *int    `json:"room"`
	Email          *string `json:"email"`
	Login          *string `json:"login"`
	Password       *string `json:"password"`
	Role           *string `json:"role"`
	RefreshToken   *string `json:"refreshToken"`
}

func (d *Doctor) DecodeFromJSON(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(d); err != nil {
		return err
	}
	return nil
}
func (db *DB) GetDoctors(offset, limit int, filter string) ([]Doctor, error) {
	var order string

	if filter != "" {
		order = fmt.Sprintf("ORDER BY %s", filter)
	}

	stmt := fmt.Sprintf(`SELECT * FROM doctors %s
    LIMIT $1 OFFSET $2;`, order)

	rows, err := db.DB.Query(stmt, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	doctors := make([]Doctor, 0, 0)

	for rows.Next() {
		doctor := Doctor{}

		err := rows.Scan(&doctor.Id, &doctor.FirstName, &doctor.LastName, &doctor.Specialization,
			&doctor.Room, &doctor.Email)

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
	stmt := `SELECT * FROM doctors WHERE id=$1 LIMIT 1;`

	row := db.DB.QueryRow(stmt, id)

	d := Doctor{}

	err := row.Scan(&d.Id, &d.FirstName, &d.LastName, &d.Specialization, &d.Room, &d.Email)

	if err != nil {
		return Doctor{}, err
	}

	return d, nil
}

func (db *DB) InsertDoctor(d Doctor) error {
	stmt := `INSERT INTO doctors (firstname, lastname, specialization, email, room) VALUES ($1, $2, $3, $4, $5);`

	if _, err := db.DB.Exec(stmt, d.FirstName, d.LastName, d.Specialization, d.Email, d.Room); err != nil {
		return err
	}

	return nil
}

//func (db *DB) UpdateDoctorAll(id int, d Doctor) error {
//	if _, err := db.GetDoctorId(id); err != nil {
//		return err
//	}
//
//	stmt := "UPDATE doctors SET  firstname = $1," +
//		"lastname = $2, specialization= $3, email = $4, room = $5 WHERE id=$6;"
//
//	if _, err := db.DB.Exec(stmt, d.FirstName, d.LastName, d.Specialization, d.Email, d.Room, id); err != nil {
//		return err
//	}
//
//	return nil
//}

func (db *DB) UpdateDoctor(id int, req string) error {
	if _, err := db.GetDoctorId(id); err != nil {
		return err
	}

	stmt := fmt.Sprintf("UPDATE doctors SET %s WHERE id=$1;", req)

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteDoctor(id int) error {
	if _, err := db.GetDoctorId(id); err != nil {
		return err
	}

	stmt := `DELETE FROM doctors WHERE id=$1;`

	if _, err := db.DB.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) CheckEmailDoctor(email string) error {
	stmt := `SELECT * FROM doctors WHERE email=$1 LIMIT 1;`

	row := db.DB.QueryRow(stmt, email)

	var doctorEmail string
	err := row.Scan(&doctorEmail)

	return err
}

func (db *DB) CheckLoginAndPasswordDoctor(login, password string) (int, error) {
	stmt := `SELECT id FROM doctors WHERE login=$1 AND password==&2 LIMIT 1;`

	row := db.DB.QueryRow(stmt, login, password)

	var id int
	err := row.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *DB) UpdateDoctorRefreshToken(id int, refreshToken string) error {
	stmt := "UPDATE patients  SET refresh_token=$1 WHERE id=$2;"

	if _, err := db.DB.Exec(stmt, refreshToken, id); err != nil {
		return err
	}

	return nil
}

func (db *DB) CheckLoginDoctor(login string) error {
	stmt := `SELECT id FROM doctors WHERE login=$1 LIMIT 1;`

	row := db.DB.QueryRow(stmt, login)

	var id int
	err := row.Scan(&id)

	return err
}
