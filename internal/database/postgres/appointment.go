package postgres

import (
	"fmt"
	"time"
)

func (db *DB) GetAppointments(offset, limit int) ([]Appointment, error) {
	stmt := fmt.Sprintf(`SELECT doctors.firstname, doctors.lastname, doctors.specialization,
      records.time_start
      FROM doctors
      JOIN records
       ON doctors.id = records.doctorid
      WHERE records.patientid IS NULL AND records.time_start > now()
      LIMIT $1 OFFSET $2;`)

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
        AND records.time_start > now()
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

//
//func (db *DB) ValidateRecord(id int, record Record) error {
//
//	validTime := record.Time_End.After(record.Time_Start) && record.Time_Start.After(time.Now())
//
//	if !validTime {
//		return recordserr.TimeErr
//	}
//
//	_, err := db.GetDoctorId(record.DoctorId)
//
//	if err != nil {
//
//		if errors.Is(err, sql.ErrNoRows) {
//			return recordserr.DoctorIdErr
//		}
//		return err
//	}
//
//	stmt := `SELECT time_start FROM records WHERE doctorid = $1 AND time_start = $2`
//
//	row := db.DB.QueryRow(stmt, record.DoctorId, record.Time_Start)
//
//	var timeStart time.Time
//
//	err = row.Scan(&timeStart)
//
//	if err != nil && !errors.Is(err, sql.ErrNoRows) {
//		return err
//	}
//
//	return nil
//
//}

func (db *DB) InsertAppointment(record Record) error {
	stmt := `INSERT INTO records (doctorid,patientid, time_start, time_end) VALUES($1, $2, $3, $4);`

	_, err := db.DB.Exec(stmt, record.DoctorId, record.PatientId, record.Time_Start, record.Time_End)

	if err != nil {
		return err
	}

	return nil

}

//
//func (db *DB) UpdateAppointmentAll(id int, start time.Time, record Record) error {
//
//	stmt := `UPDATE records SET doctorid = $1,
//                   patientid = $2,
//                   time_start = $3,
//                   time_end = $4
//                   WHERE doctorid=$5 AND time_start=$6`
//
//	_, err := db.DB.Exec(stmt, record.DoctorId, record.PatientId, record.Time_Start, record.Time_End,
//		id, start)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//
//}

func (db *DB) UpdateAppointment(id int, start time.Time, req string) error {
	stmt := fmt.Sprintf("UPDATE records SET %s WHERE doctorid=$1 AND time_start=$2;", req)

	_, err := db.DB.Exec(stmt, id, start)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteAppointment(id int, start time.Time) error {

	if err := db.CheckRecord(id, start); err != nil {
		return err
	}

	stmt := `DELETE FROM records WHERE doctorid=$1 AND time_start=$2;`

	if _, err := db.DB.Exec(stmt, id, start); err != nil {
		return err
	}

	return nil
}

func (db *DB) CheckRecord(doctorId int, timeStart time.Time) error {
	stmt := `SELECT * FROM records WHERE doctorid=$1 AND time_start=$2 LIMIT 1;`

	row := db.DB.QueryRow(stmt, doctorId, timeStart)

	var record Record

	err := row.Scan(&record.DoctorId, &record.PatientId, &record.Time_End, &record.Time_End)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetAppointmentsCountAllDoctors(offset, limit int, filterId string) ([]AppointmentsCount, error) {
	var req string

	if len(filterId) > 0 {
		req = fmt.Sprintf("AND doctors.id IN (%s)", filterId)
	}

	stmt := fmt.Sprintf(`SELECT doctors.id, doctors.firstname, doctors.lastname, doctors.specialization,
      count(records.patientid) AS Appointments
      FROM doctors
      JOIN records
       ON doctors.id = records.doctorid
      WHERE records.time_start > now() %s
      GROUP BY doctors.id
      ORDER BY Appointments
      LIMIT $1 OFFSET $2;`, req)

	rows, err := db.DB.Query(stmt, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appointmentsCount := make([]AppointmentsCount, 0, 0)

	for rows.Next() {
		appointmentCount := AppointmentsCount{}

		err := rows.Scan(&appointmentCount.DoctorId, &appointmentCount.DoctorFirstName,
			&appointmentCount.DoctorLastName, &appointmentCount.DoctorSpecialization, &appointmentCount.Appointments)

		if err != nil {
			return nil, err
		}

		appointmentsCount = append(appointmentsCount, appointmentCount)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appointmentsCount, nil
}
