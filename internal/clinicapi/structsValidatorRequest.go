package clinicapi

import (
	"database/sql"
	"errors"
	"github.com/Vadim992/clinicAPI/internal/customErr/recordsErr"
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"github.com/Vadim992/clinicAPI/pkg/validator/validate"
	"net/http"
	"time"
)

// validateNumPage checks that page > 0
func validateNumPage(page, pageSize int) bool {
	if page < 1 || pageSize < 1 {
		return false
	}
	return true
}

//validateEmail checks Is valid PATIENT's email or not

func validateEmail(email sql.NullString) bool {
	if email.Valid {
		return validate.ValidateEmail(email.String)
	}

	return true
}

// func validatePhoneNum checks Is valid PATIENT's phone number or not
func validatePhoneNum(phoneNum string) bool {
	return validate.ValidatePhoneNum(phoneNum)
}

func (c *ClinicAPI) validatePatientEmailPhone(w http.ResponseWriter, r *http.Request, patient *postgres.Patient) bool {

	if !validateEmail(patient.Email) {
		c.clientErr(w, http.StatusBadRequest)
		return false
	}

	if patient.Email.Valid {

		err := c.DB.CheckEmailPatient(patient.Email.String)

		if !errors.Is(err, sql.ErrNoRows) {
			c.clientErr(w, http.StatusBadRequest)
			return false
		}

	}

	switch r.Method {

	case http.MethodPatch:

		if patient.Phone_Number != "" && !validatePhoneNum(patient.Phone_Number) {
			c.clientErr(w, http.StatusBadRequest)
			return false
		}
	default:
		if !validatePhoneNum(patient.Phone_Number) {
			c.clientErr(w, http.StatusBadRequest)
			return false
		}

	}

	if patient.Phone_Number != "" {

		patient.Phone_Number = formatPhoneNumber(patient.Phone_Number)
		err := c.DB.CheckPhoneNumberPatient(patient.Phone_Number)

		if !errors.Is(err, sql.ErrNoRows) {
			c.clientErr(w, http.StatusBadRequest)
			return false
		}
	}

	return true

}

func (c *ClinicAPI) validateDoctor(w http.ResponseWriter, r *http.Request, doctor postgres.Doctor) bool {

	switch r.Method {

	case http.MethodPatch:

		if doctor.Email != "" {

			if !validate.ValidateEmail(doctor.Email) {
				c.clientErr(w, http.StatusBadRequest)
				return false
			}

			err := c.DB.CheckEmailDoctor(doctor.Email)

			if !errors.Is(err, sql.ErrNoRows) {
				c.clientErr(w, http.StatusBadRequest)
				return false
			}

		}

	default:
		if !putCheckStructs(doctor) {
			c.clientErr(w, http.StatusBadRequest)
			return false
		}

		if !validate.ValidateEmail(doctor.Email) {
			c.clientErr(w, http.StatusBadRequest)
			return false
		}

		err := c.DB.CheckEmailDoctor(doctor.Email)
		if !errors.Is(err, sql.ErrNoRows) {
			c.clientErr(w, http.StatusBadRequest)
			return false
		}

	}

	return true

}

func validateDoctorPatientIDs(DoctorId, PatientId int) error {

	if !(validate.ValidateId(DoctorId) && validate.ValidateId(PatientId)) {
		return recordsErr.IdErr
	}
	return nil
}

func (c *ClinicAPI) validateRecordPutPost(record postgres.Record) error {
	//check that time_start and time_end are ok

	validTime := record.Time_End.After(record.Time_Start) && record.Time_Start.After(time.Now())

	if !validTime {
		return recordsErr.TimeErr
	}

	//check that doctorId and patientId are ok ( id > 0)

	if record.PatientId.Valid {

		doctorId := record.DoctorId
		patientId := int(record.PatientId.Int64)

		if err := validateDoctorPatientIDs(doctorId, patientId); err != nil {
			return recordsErr.IdErr
		}

		//check that patientId exists
		_, err := c.DB.GetPatientId(int(record.PatientId.Int64))

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				return recordsErr.PatientIdErr
			}

			return err
		}
	} else {

		if !validate.ValidateId(record.DoctorId) {
			return recordsErr.IdErr
		}
	}

	//check that doctorId exists
	_, err := c.DB.GetDoctorId(record.DoctorId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return recordsErr.DoctorIdErr
		}
		return err
	}
	return nil
}

func (c *ClinicAPI) validateRecordPOST(record postgres.Record) error {

	err := c.validateRecordPutPost(record)

	if err != nil {
		return err
	}
	err = c.DB.CheckRecord(record.DoctorId, record.Time_Start)

	if !errors.Is(err, sql.ErrNoRows) {
		if err == nil {
			return recordsErr.InsertErr
		}
		return err
	}

	return nil

}

func (c *ClinicAPI) validateRecordPUT(id int, data AppointmentsData) error {

	startTime := data.StartTime

	// check that this appointment exists
	err := c.DB.CheckRecord(id, startTime)

	if err != nil {

		return err
	}

	record := data.Record

	err = c.validateRecordPutPost(record)

	if err != nil {
		return err
	}

	return nil

}

func (c *ClinicAPI) validateRecordPATCH(id int, data AppointmentsData) error {

	startTime := data.StartTime

	// check that this appointment exists
	err := c.DB.CheckRecord(id, startTime)

	if err != nil {
		return err
	}

	record := data.Record

	// check that time_start and time_end is ok
	if err := c.validateRecordTime(id, startTime, record); err != nil {
		return err
	}

	// check that patientId is ok
	if record.PatientId.Valid {

		if !validate.ValidateId(int(record.PatientId.Int64)) {
			return recordsErr.IdErr
		}

		//check that patientId exists
		_, err := c.DB.GetPatientId(int(record.PatientId.Int64))

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				return recordsErr.PatientIdErr
			}

			return err
		}
	}
	// check that doctorId is ok
	if record.DoctorId != 0 {

		if !validate.ValidateId(record.DoctorId) {
			return recordsErr.IdErr
		}

		//check that doctorId exists
		_, err := c.DB.GetDoctorId(record.DoctorId)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				return recordsErr.DoctorIdErr
			}
			return err
		}
	}

	return nil

}

func (c *ClinicAPI) validateRecordTime(id int, start time.Time, record postgres.Record) error {

	switch true {
	case !record.Time_Start.IsZero() && !record.Time_End.IsZero():

		validTime := record.Time_End.After(record.Time_Start) && record.Time_Start.After(time.Now())

		if !validTime {
			return recordsErr.TimeErr
		}

	case !record.Time_Start.IsZero():

		var timeEnd time.Time

		stmt := `SELECT time_end FROM records
                 WHERE doctorid = $1 AND time_start = $2`

		row := c.DB.DB.QueryRow(stmt, id, start)

		err := row.Scan(&timeEnd)

		if err != nil {
			return err
		}

		validTime := record.Time_Start.Before(timeEnd) && record.Time_Start.After(time.Now())

		if !validTime {
			return recordsErr.TimeErr
		}

	case !record.Time_End.IsZero():

		validTime := record.Time_End.After(start) && record.Time_End.After(time.Now())

		if !validTime {
			return recordsErr.TimeErr
		}

	}

	return nil

}
