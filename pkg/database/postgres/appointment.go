package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Vadim992/clinicAPI/internal/customErr/recordsErr"
	"time"
)

type Record struct {
	DoctorId  int           `json:"doctorId"`
	PatientId sql.NullInt64 `json:"patientId"`
	TimeStart time.Time     `json:"timeStart"`
	TimeEnd   time.Time     `json:"timeEnd"`
}

type Appointment struct {
	DoctorFirstName      string    `json:"doctorFirstName"`
	DoctorLastName       string    `json:"doctorLastName"`
	DoctorSpecialization string    `json:"doctorSpecialization"`
	TimeStart            time.Time `json:"timeStart"`
}

func (db *DB) GetAppointments(offset, limit int, filter string) ([]Appointment, error) {

	var order string

	if filter != "" {
		order = fmt.Sprintf("ORDER BY %s", filter)
	}

	stmt := fmt.Sprintf(`SELECT doctors.firstname, doctors.lastname, doctors.specialization,
       records.time_start
       FROM doctors
       JOIN records
        ON doctors.id = records.doctorid
       WHERE records.patientid IS NULL AND records.time_start > current_timestamp(0)
       %s
       LIMIT $1 OFFSET $2;`, order)

	rows, err := db.DB.Query(stmt, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appointments := make([]Appointment, 0, 0)

	for rows.Next() {
		appointment := Appointment{}

		err := rows.Scan(&appointment.DoctorFirstName, &appointment.DoctorLastName,
			&appointment.DoctorSpecialization, &appointment.TimeStart)

		if err != nil {
			return nil, err
		}

		appointments = append(appointments, appointment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil

}

func (db *DB) GetAppointmentsId(offset, limit, id int) ([]Appointment, error) {
	_, err := db.GetDoctorId(id)

	if err != nil {
		return nil, err
	}

	stmt := `SELECT doctors.firstname, doctors.lastname, doctors.specialization,
       records.time_start
       FROM doctors
       JOIN records
        ON doctors.id = records.doctorid
       WHERE doctors.id = $1 AND 
             records.patientid IS NULL 
         AND records.time_start > current_timestamp(0)
       LIMIT $2 OFFSET $3;`

	rows, err := db.DB.Query(stmt, id, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appointments := make([]Appointment, 0, 0)

	for rows.Next() {
		appointment := Appointment{}

		err := rows.Scan(&appointment.DoctorFirstName, &appointment.DoctorLastName,
			&appointment.DoctorSpecialization, &appointment.TimeStart)

		if err != nil {
			return nil, err
		}

		appointments = append(appointments, appointment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil

}

func (db *DB) ValidateRecord(record Record) error {

	validTime := record.TimeEnd.After(record.TimeStart) && record.TimeStart.After(time.Now())

	if !validTime {
		return recordsErr.TimeErr
	}

	//_, err := db.GetDoctorId(id)
	//
	//if err != nil {
	//
	//	if !errors.Is(err, sql.ErrNoRows) {
	//		return err
	//	}
	//}

	_, err := db.GetDoctorId(record.DoctorId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return recordsErr.DoctorIdErr
		}
		return err
	}

	stmt := `SELECT time_start FROM records WHERE doctorid = $1 AND time_start = $2`

	row := db.DB.QueryRow(stmt, record.DoctorId, record.TimeStart)

	var timeStart time.Time

	err = row.Scan(&timeStart)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	return nil

}

func (db *DB) InsertAppointment(record Record) error {

	stmt := `INSERT INTO records (doctorid,patientid, time_start, time_end) VALUES($1, $2, $3, $4)`

	_, err := db.DB.Exec(stmt, record.DoctorId, record.PatientId, record.TimeStart, record.TimeEnd)

	if err != nil {
		return err
	}

	return nil

}

func (db *DB) UpdateAppointmentAll(id int, start time.Time, record Record) error {

	_, err := db.GetDoctorId(id)

	if err != nil {
		return err
	}

	stmt := `UPDATE records SET doctorid = $1,
                   patientid = $2,
                   time_start = $3, 
                   time_end = $4`

	_, err = db.DB.Exec(stmt, record.DoctorId, record.PatientId, record.TimeStart, record.TimeEnd)

	if err != nil {
		return err
	}

	return nil

}

func (db *DB) UpdateAppointment(id int, start time.Time, req string) error {

	_, err := db.GetDoctorId(id)

	if err != nil {
		return err
	}

	stmt := fmt.Sprintf("UPDATE recotrds SET %s WHERE id=$1 AND time_start=$2", req)

	_, err = db.DB.Exec(stmt, id, start)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteAppointment(id int, start time.Time) error {

	_, err := db.GetDoctorId(id)

	if err != nil {
		return err
	}

	stmt := `DELETE FROM records WHERE id=$1 AND time_start=$2`

	if _, err := db.DB.Exec(stmt, id, start); err != nil {
		return err
	}

	return nil
	return nil
}
