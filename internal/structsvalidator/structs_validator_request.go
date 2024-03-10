package structsvalidator

import (
	"database/sql"
	"errors"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"github.com/Vadim992/clinicAPI/internal/dto"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/patientserr"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/recordserr"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/structserr"
	"github.com/Vadim992/clinicAPI/internal/helpers/dbhelpers"
	"github.com/Vadim992/clinicAPI/pkg/validator/validate"
	"strings"
	"time"
)

// validateNumPage checks that page > 0
func ValidateNumPage(page, pageSize int) bool {
	if page < 1 || pageSize < 1 {
		return false
	}

	return true
}

func ValidatePatientEmailPhone(patient postgres.Patient) error {
	if patient.Email != nil {
		*patient.Email = strings.ToLower(*patient.Email)

		if !validate.ValidateEmail(*patient.Email) {
			return structserr.EmailErr
		}

		err := postgres.DataBase.CheckEmailPatient(*patient.Email)

		if err == nil {
			return recordserr.InsertErr
		}

		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	if patient.Phone_Number != nil {
		if !validate.ValidatePhoneNum(*patient.Phone_Number) {
			return patientserr.PhoneErr
		}

		*patient.Phone_Number = dbhelpers.FormatPhoneNumber(*patient.Phone_Number)

		err := postgres.DataBase.CheckPhoneNumberPatient(*patient.Phone_Number)

		if err == nil {
			return recordserr.InsertErr
		}

		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	return nil
}

func ValidateDoctorEmail(doctor postgres.Doctor) error {
	if doctor.Email != nil {
		*doctor.Email = strings.ToLower(*doctor.Email)

		if !validate.ValidateEmail(*doctor.Email) {
			return structserr.EmailErr
		}

		err := postgres.DataBase.CheckEmailDoctor(*doctor.Email)

		if err == nil {
			return recordserr.InsertErr
		}

		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	return nil
}

func ValidateRecordPutPost(record postgres.Record) error {
	//check that time_start and time_end are ok
	validTime := record.Time_End.After(*record.Time_Start) && record.Time_Start.After(time.Now())

	if !validTime {
		return recordserr.TimeErr
	}

	//check that doctorId and patientId are ok ( id > 0)
	if record.PatientId != nil {
		patientId := *record.PatientId

		if !validate.ValidateId(patientId) {
			return recordserr.IdErr
		}

		//check that patientId exists
		_, err := postgres.DataBase.GetPatientId(patientId)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return recordserr.PatientIdErr
			}

			return err
		}
	}
	if !validate.ValidateId(*record.DoctorId) {
		return recordserr.IdErr
	}

	//check that doctorId exists
	_, err := postgres.DataBase.GetDoctorId(*record.DoctorId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return recordserr.DoctorIdErr
		}

		return err
	}

	return nil
}

func ValidateRecordPOST(record postgres.Record) error {
	err := ValidateRecordPutPost(record)

	if err != nil {
		return err
	}

	err = postgres.DataBase.CheckRecord(*record.DoctorId, *record.Time_Start)

	if !errors.Is(err, sql.ErrNoRows) {
		if err == nil {
			return recordserr.InsertErr
		}

		return err
	}

	return nil
}

//func (c *ClinicAPI) validateRecordPUT(id int, data AppointmentsData) error {
//
//	startTime := data.StartTime
//
//	// check that this appointment exists
//	err := c.DB.CheckRecord(id, startTime)
//
//	if err != nil {
//
//		return err
//	}
//
//	record := data.Record
//
//	err = c.validateRecordPutPost(record)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//
//}

func ValidateRecordPATCH(id int, data dto.AppointmentsData) error {
	startTime := data.StartTime

	// check that this appointment exists
	err := postgres.DataBase.CheckRecord(id, startTime)

	if err != nil {
		return err
	}

	record := data.Record

	// check that time_start and time_end is ok
	if err := validateRecordTime(id, startTime, record); err != nil {
		return err
	}

	// check that patientId is ok
	if record.PatientId != nil {
		if !validate.ValidateId(*record.PatientId) {
			return recordserr.IdErr
		}

		//check that patientId exists
		_, err := postgres.DataBase.GetPatientId(*record.PatientId)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return recordserr.PatientIdErr
			}

			return err
		}
	}
	// check that doctorId is ok
	if record.DoctorId != nil {
		if !validate.ValidateId(*record.DoctorId) {
			return recordserr.IdErr
		}

		//check that doctorId exists
		_, err := postgres.DataBase.GetDoctorId(*record.DoctorId)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return recordserr.DoctorIdErr
			}

			return err
		}
	}

	return nil

}

func validateRecordTime(id int, start time.Time, record postgres.Record) error {
	switch {
	// when Time_Start and Time_End are exist
	case record.Time_Start != nil && record.Time_End != nil:
		validTime := record.Time_End.After(*record.Time_Start) && record.Time_Start.After(time.Now())

		if !validTime {
			return recordserr.TimeErr
		}
		// when only Time_Start  exists
	case record.Time_Start != nil:
		var timeEnd time.Time

		stmt := `SELECT time_end FROM records
                WHERE doctorid = $1 AND time_start = $2;`

		row := postgres.DataBase.DB.QueryRow(stmt, id, start)

		err := row.Scan(&timeEnd)

		if err != nil {
			return err
		}

		validTime := record.Time_Start.Before(timeEnd) && record.Time_Start.After(time.Now())

		if !validTime {
			return recordserr.TimeErr
		}
	// when only Time_End  exists
	case record.Time_End != nil:
		validTime := record.Time_End.After(start) && record.Time_End.After(time.Now())

		if !validTime {
			return recordserr.TimeErr
		}
	}

	return nil
}
