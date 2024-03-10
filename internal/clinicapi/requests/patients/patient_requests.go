package patients

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"github.com/Vadim992/clinicAPI/internal/dto"
	"github.com/Vadim992/clinicAPI/internal/helpers"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/dtoerr"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/patientserr"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/recordserr"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/structserr"
	"github.com/Vadim992/clinicAPI/internal/helpers/dbhelpers"
	"github.com/Vadim992/clinicAPI/internal/structsvalidator"
	"github.com/Vadim992/clinicAPI/pkg/logger"
	"github.com/Vadim992/clinicAPI/pkg/validator/validate"
	"io"
	"net/http"
)

func GetPatients(w http.ResponseWriter, r *http.Request) {
	err := getPatients(w, r)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, io.EOF):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, dtoerr.PageErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		default:
			helpers.ServeErr(w, err)
		}
	}
}

func getPatients(w http.ResponseWriter, r *http.Request) error {
	var pageData dto.PageData

	if err := pageData.DecodeFromJSON(r); err != nil {
		return err
	}

	if !structsvalidator.ValidateNumPage(pageData.Page, pageData.PageSize) {
		return dtoerr.PageErr
	}

	offset := (pageData.Page - 1) * pageData.PageSize

	var filter string

	switch {
	case pageData.PatientsFilter.PhoneFilter && pageData.PatientsFilter.FirstNameFilter:
		filter = "phone_number, firstname"
	case pageData.PatientsFilter.PhoneFilter:
		filter = "phone_number"
	case pageData.PatientsFilter.FirstNameFilter:
		filter = "firstname"
	}

	patients, err := postgres.DataBase.GetPatients(offset, pageData.PageSize, filter)

	if err != nil {
		return err
	}

	b, err := json.Marshal(patients)

	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func GetPatientId(w http.ResponseWriter, r *http.Request) {
	id, err := dbhelpers.ConvertIdFromStrToInt(r)

	if err != nil {
		helpers.ServeErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		helpers.NotFound(w)
		return
	}

	err = getPatientId(w, id)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			helpers.NotFound(w)
		default:
			helpers.ServeErr(w, err)
		}
	}
}

func getPatientId(w http.ResponseWriter, id int) error {
	patient, err := postgres.DataBase.GetPatientId(id)

	if err != nil {
		return err
	}

	b, err := json.Marshal(patient)

	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func PostPatient(w http.ResponseWriter, r *http.Request) {
	err := postPatient(w, r)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, io.EOF):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.InsertIdErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, structserr.EmailErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, patientserr.PhoneErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.InsertErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, structserr.EmptyFieldErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, structserr.InvalidTypeOfStructErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		default:
			helpers.ServeErr(w, err)
		}
	}
}

func postPatient(w http.ResponseWriter, r *http.Request) error {
	var patient postgres.Patient

	if err := patient.DecodeFromJSON(r); err != nil {
		return err
	}

	if patient.Id != nil {
		return recordserr.InsertIdErr
	}

	err := structsvalidator.ValidatePatientEmailPhone(patient)

	if err != nil {
		return err
	}

	err = dbhelpers.CheckStructsFields(patient)

	if err != nil {
		return err
	}

	if err := postgres.DataBase.InsertPatient(patient); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

// func (c *ClinicAPI) putPatient(w http.ResponseWriter, r *http.Request, id int) {
//
//	var patient postgres.Patient
//
//	if err := clinicapi.decode(r, &patient); err != nil {
//		if errors.Is(err, io.EOF) {
//			c.clientErr(w, http.StatusBadRequest)
//			return
//		}
//		c.serveErr(w, err)
//		return
//	}
//
//	patient.Email.String = strings.ToLower(patient.Email.String)
//
//	if !c.validatePatientEmailPhone(w, r, &patient) {
//		return
//	}
//
//	//patient.Phone_Number = formatPhoneNumber(patient.Phone_Number)
//
//	if !clinicapi.putCheckStructs(patient) {
//		c.clientErr(w, http.StatusBadRequest)
//		return
//	}
//
//	if err := c.DB.UpdatePatientAll(id, patient); err != nil {
//
//		if errors.Is(err, sql.ErrNoRows) {
//			c.notFound(w)
//			return
//		}
//
//		c.serveErr(w, err)
//		return
//	}
//
// }
func PatchPatientId(w http.ResponseWriter, r *http.Request) {
	id, err := dbhelpers.ConvertIdFromStrToInt(r)

	if err != nil {
		helpers.ServeErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		helpers.NotFound(w)
		return
	}

	err = patchPatientId(r, id)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, io.EOF):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.InsertIdErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, structserr.EmailErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, patientserr.PhoneErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.InsertErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, structserr.InvalidTypeOfStructErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, structserr.NoDataForPatchReqErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, sql.ErrNoRows):
			helpers.NotFound(w)
		default:
			helpers.ServeErr(w, err)
		}
	}
}

func patchPatientId(r *http.Request, id int) error {
	var patient postgres.Patient

	if err := patient.DecodeFromJSON(r); err != nil {
		return err
	}

	if patient.Id != nil {
		return recordserr.InsertIdErr
	}

	err := structsvalidator.ValidatePatientEmailPhone(patient)

	if err != nil {
		return err
	}

	req, err := dbhelpers.CheckStructFieldsPatch(patient)

	if err != nil {
		return err
	}

	if req == "" {
		return structserr.NoDataForPatchReqErr
	}

	if err := postgres.DataBase.UpdatePatient(id, req); err != nil {
		return err
	}

	return nil
}

func DeletePatientId(w http.ResponseWriter, r *http.Request) {
	id, err := dbhelpers.ConvertIdFromStrToInt(r)

	if err != nil {
		helpers.ServeErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		helpers.NotFound(w)
		return
	}

	err = deletePatientId(id)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			helpers.NotFound(w)
		default:
			helpers.ServeErr(w, err)
		}
	}
}

func deletePatientId(id int) error {
	if err := postgres.DataBase.DeletePatient(id); err != nil {
		return err
	}

	return nil
}
