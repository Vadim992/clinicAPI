package appointments

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	"strings"
)

func GetAppointments(w http.ResponseWriter, r *http.Request) {
	err := getAppointments(w, r)

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

func getAppointments(w http.ResponseWriter, r *http.Request) error {
	var pageData dto.PageData

	if err := pageData.DecodeFromJSON(r); err != nil {
		return err
	}

	if !structsvalidator.ValidateNumPage(pageData.Page, pageData.PageSize) {
		return dtoerr.PageErr
	}

	offset := (pageData.Page - 1) * pageData.PageSize

	appointments, err := postgres.DataBase.GetAppointments(offset, pageData.PageSize)

	if err != nil {
		return err
	}

	b, err := json.Marshal(appointments)

	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func GetAppointmentId(w http.ResponseWriter, r *http.Request) {
	id, err := dbhelpers.ConvertIdFromStrToInt(r)

	if err != nil {
		helpers.ServeErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		helpers.NotFound(w)
		return
	}

	err = getAppointmentId(w, r, id)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, io.EOF):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, dtoerr.PageErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, sql.ErrNoRows):
			helpers.NotFound(w)
		default:
			helpers.ServeErr(w, err)

		}
	}

}

func getAppointmentId(w http.ResponseWriter, r *http.Request, id int) error {
	var pageData dto.PageData

	if err := pageData.DecodeFromJSON(r); err != nil {
		return err
	}

	if !structsvalidator.ValidateNumPage(pageData.Page, pageData.PageSize) {
		return dtoerr.PageErr
	}

	offset := (pageData.Page - 1) * pageData.PageSize

	appointments, err := postgres.DataBase.GetAppointmentsId(offset, pageData.PageSize, id)

	if err != nil {
		return err
	}

	b, err := json.Marshal(appointments)

	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func PostAppointment(w http.ResponseWriter, r *http.Request) {
	err := postAppointment(w, r)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, io.EOF):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, structserr.InvalidTypeOfStructErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, structserr.EmptyFieldErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.TimeErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.IdErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.PatientIdErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.DoctorIdErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.InsertErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		default:
			helpers.ServeErr(w, err)
		}
	}
}

func postAppointment(w http.ResponseWriter, r *http.Request) error {
	var record postgres.Record

	if err := record.DecodeFromJSON(r); err != nil {
		return err
	}

	if err := dbhelpers.CheckStructsFields(record); err != nil {
		return err
	}

	if err := structsvalidator.ValidateRecordPOST(record); err != nil {
		return err
	}

	if err := postgres.DataBase.InsertAppointment(record); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

//
//func (c *ClinicAPI) putAppointment(w http.ResponseWriter, r *http.Request, id int) {
//
//	var appointmentData clinicapi.AppointmentsData
//
//	if err := clinicapi.decode(r, &appointmentData); err != nil {
//		if errors.Is(err, io.EOF) {
//			c.clientErr(w, http.StatusBadRequest)
//			return
//		}
//		c.serveErr(w, err)
//		return
//	}
//
//	if !clinicapi.putCheckStructs(appointmentData.Record) {
//		c.clientErr(w, http.StatusBadRequest)
//		return
//	}
//
//	if err := c.validateRecordPUT(id, appointmentData); err != nil {
//
//		switch true {
//		case errors.Is(err, recordsErr.IdErr):
//
//			c.clientErr(w, http.StatusBadRequest)
//			return
//
//		case errors.Is(err, recordsErr.TimeErr):
//
//			c.clientErr(w, http.StatusBadRequest)
//			return
//
//		case errors.Is(err, recordsErr.DoctorIdErr):
//
//			c.clientErr(w, http.StatusBadRequest)
//			return
//
//		case errors.Is(err, recordsErr.PatientIdErr):
//
//			c.clientErr(w, http.StatusBadRequest)
//			return
//
//		case errors.Is(err, sql.ErrNoRows):
//
//			c.clientErr(w, http.StatusBadRequest)
//			return
//
//		default:
//			c.serveErr(w, err)
//			return
//		}
//	}
//
//	start := appointmentData.StartTime
//	record := appointmentData.Record
//
//	if err := c.DB.UpdateAppointmentAll(id, start, record); err != nil {
//		c.serveErr(w, err)
//		return
//	}
//
//}

func PatchAppointmentId(w http.ResponseWriter, r *http.Request) {
	id, err := dbhelpers.ConvertIdFromStrToInt(r)

	if err != nil {
		helpers.ServeErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		helpers.NotFound(w)
		return
	}

	err = patchAppointmentId(r, id)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, io.EOF):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, sql.ErrNoRows):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.TimeErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.IdErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.PatientIdErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.DoctorIdErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, structserr.NoDataForPatchReqErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		default:
			helpers.ServeErr(w, err)
		}
	}
}

func patchAppointmentId(r *http.Request, id int) error {
	var appointmentData dto.AppointmentsData

	if err := appointmentData.DecodeFromJSON(r); err != nil {
		return err
	}

	if err := structsvalidator.ValidateRecordPATCH(id, appointmentData); err != nil {
		return err
	}

	record := appointmentData.Record
	req, err := dbhelpers.CheckStructFieldsPatch(record)

	if err != nil {
		return err
	}

	if req == "" {
		return structserr.NoDataForPatchReqErr
	}

	start := appointmentData.StartTime

	if err := postgres.DataBase.UpdateAppointment(id, start, req); err != nil {
		return err
	}

	return nil
}

func DeleteAppointmentId(w http.ResponseWriter, r *http.Request) {
	id, err := dbhelpers.ConvertIdFromStrToInt(r)

	if err != nil {
		logger.ErrLog.Println(err)
		helpers.ServeErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		helpers.NotFound(w)
		return
	}

	err = deleteAppointmentId(r, id)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, io.EOF):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, sql.ErrNoRows):
			helpers.ClientErr(w, http.StatusBadRequest)
		default:
			helpers.ServeErr(w, err)
		}
	}
}

func deleteAppointmentId(r *http.Request, id int) error {
	var appointmentData dto.AppointmentsData

	if err := appointmentData.DecodeFromJSON(r); err != nil {
		return err
	}

	if err := postgres.DataBase.DeleteAppointment(id, appointmentData.StartTime); err != nil {
		return err
	}

	return nil
}

func GetAppointmentsCountAllDoctors(w http.ResponseWriter, r *http.Request) {
	err := getAppointmentsCountAllDoctors(w, r)

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

func getAppointmentsCountAllDoctors(w http.ResponseWriter, r *http.Request) error {
	var pageData dto.PageData

	if err := pageData.DecodeFromJSON(r); err != nil {
		return err
	}

	if !structsvalidator.ValidateNumPage(pageData.Page, pageData.PageSize) {
		return dtoerr.PageErr
	}

	offset := (pageData.Page - 1) * pageData.PageSize

	filterId := convertIntSliceToString(pageData.AppointmentsDoctorId)

	appointmentsCount, err := postgres.DataBase.GetAppointmentsCountAllDoctors(offset, pageData.PageSize, filterId)

	if err != nil {
		return err
	}

	b, err := json.Marshal(appointmentsCount)

	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func convertIntSliceToString(ids []int) string {
	if ids == nil || len(ids) == 0 {
		return ""
	}
	var b strings.Builder

	for _, val := range ids {
		b.WriteString(fmt.Sprintf("%d, ", val))
	}

	str := strings.TrimSpace(b.String())

	return str[:len(str)-1]
}
