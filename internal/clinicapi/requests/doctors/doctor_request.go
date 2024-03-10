package doctors

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"github.com/Vadim992/clinicAPI/internal/dto"
	"github.com/Vadim992/clinicAPI/internal/helpers"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/dtoerr"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/recordserr"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/structserr"
	"github.com/Vadim992/clinicAPI/internal/helpers/dbhelpers"
	"github.com/Vadim992/clinicAPI/internal/structsvalidator"
	"github.com/Vadim992/clinicAPI/pkg/logger"
	"github.com/Vadim992/clinicAPI/pkg/validator/validate"
	"io"
	"net/http"
)

func GetDoctors(w http.ResponseWriter, r *http.Request) {
	err := getDoctors(w, r)

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

func getDoctors(w http.ResponseWriter, r *http.Request) error {
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
	case pageData.DoctorsFilter.SpecializationFilter && pageData.DoctorsFilter.FirstNameFilter:
		filter = "specialization, firstname"
	case pageData.DoctorsFilter.SpecializationFilter:
		filter = "specialization"
	case pageData.DoctorsFilter.FirstNameFilter:
		filter = "firstname"
	}

	doctors, err := postgres.DataBase.GetDoctors(offset, pageData.PageSize, filter)

	if err != nil {
		return err
	}

	b, err := json.Marshal(doctors)

	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func GetDoctorId(w http.ResponseWriter, r *http.Request) {
	id, err := dbhelpers.ConvertIdFromStrToInt(r)

	if err != nil {
		helpers.ServeErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		helpers.NotFound(w)
		return
	}

	err = getDoctorId(w, id)

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

func getDoctorId(w http.ResponseWriter, id int) error {
	doctor, err := postgres.DataBase.GetDoctorId(id)

	if err != nil {
		return err
	}

	b, err := json.Marshal(doctor)

	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func PostDoctor(w http.ResponseWriter, r *http.Request) {
	err := postDoctor(w, r)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, io.EOF):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.InsertIdErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.InsertErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, structserr.InvalidTypeOfStructErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		default:
			helpers.ServeErr(w, err)
		}
	}
}

func postDoctor(w http.ResponseWriter, r *http.Request) error {
	var doctor postgres.Doctor

	if err := doctor.DecodeFromJSON(r); err != nil {
		return err
	}

	if doctor.Id != nil {
		return recordserr.InsertIdErr
	}

	if err := structsvalidator.ValidateDoctorEmail(doctor); err != nil {
		return err
	}

	if err := dbhelpers.CheckStructsFields(doctor); err != nil {
		return err
	}

	if err := postgres.DataBase.InsertDoctor(doctor); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

// func (c *ClinicAPI) putDoctor(w http.ResponseWriter, r *http.Request, id int) {
//
//		var doctor postgres.Doctor
//
//		if err := clinicapi.decode(r, &doctor); err != nil {
//			if errors.Is(err, io.EOF) {
//				c.clientErr(w, http.StatusBadRequest)
//				return
//			}
//			c.serveErr(w, err)
//			return
//		}
//
//		doctor.Email = strings.ToLower(doctor.Email)
//
//		if !c.validateDoctor(w, r, doctor) {
//			return
//		}
//
//		if err := c.DB.UpdateDoctorAll(id, doctor); err != nil {
//
//			if errors.Is(err, sql.ErrNoRows) {
//				c.notFound(w)
//				return
//			}
//
//			c.serveErr(w, err)
//			return
//		}
//	}
func PatchDoctorId(w http.ResponseWriter, r *http.Request) {
	id, err := dbhelpers.ConvertIdFromStrToInt(r)

	if err != nil {
		helpers.ServeErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		helpers.NotFound(w)
		return
	}

	err = patchDoctorId(r, id)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, io.EOF):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.InsertIdErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, structserr.EmailErr):
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

func patchDoctorId(r *http.Request, id int) error {
	var doctor postgres.Doctor

	if err := doctor.DecodeFromJSON(r); err != nil {
		return err
	}

	if doctor.Id != nil {
		return recordserr.InsertIdErr
	}

	err := structsvalidator.ValidateDoctorEmail(doctor)

	if err != nil {
		return err
	}

	req, err := dbhelpers.CheckStructFieldsPatch(doctor)

	if err != nil {
		return err
	}

	if req == "" {
		return structserr.NoDataForPatchReqErr
	}

	if err := postgres.DataBase.UpdateDoctor(id, req); err != nil {
		return err
	}

	return nil
}

func DeleteDoctorId(w http.ResponseWriter, r *http.Request) {
	id, err := dbhelpers.ConvertIdFromStrToInt(r)

	if err != nil {
		helpers.ServeErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		helpers.NotFound(w)
		return
	}

	err = deleteDoctorId(id)

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

func deleteDoctorId(id int) error {
	if err := postgres.DataBase.DeleteDoctor(id); err != nil {
		return err
	}

	return nil
}
