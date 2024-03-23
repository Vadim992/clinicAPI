package clinicapi

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Vadim992/clinicAPI/internal/clinicapi/requests/patients"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"github.com/Vadim992/clinicAPI/internal/dto"
	"github.com/Vadim992/clinicAPI/internal/helpers"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/patientserr"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/recordserr"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/structserr"
	"github.com/Vadim992/clinicAPI/internal/helpers/dbhelpers"
	"github.com/Vadim992/clinicAPI/internal/jwtgen"
	"github.com/Vadim992/clinicAPI/internal/structsvalidator"
	"github.com/Vadim992/clinicAPI/pkg/logger"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

var (
	NoRoleErr = errors.New("have no 'Role' field in request body or unknown 'Role'")

	NoLoginErr    = errors.New("have no 'LogIn' field in request body")
	NoPasswordErr = errors.New("have no 'Password' field in request body")

	NoPatientStructErr = errors.New("have no 'Patient' field in request body")
	NoDoctorStructErr  = errors.New("have no 'Doctor' field in request body")

	IncorrectLoginOrPasswordErr = errors.New("incorrect login or password")

	LenLoginPasswordErr = errors.New("length of password or login is shorter than 8")
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	err := logIn(w, r)

	if err != nil {
		logger.ErrLog.Println(err)

		switch {
		case errors.Is(err, NoRoleErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, NoLoginErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, NoPasswordErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, LenLoginPasswordErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, NoPatientStructErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, NoDoctorStructErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordserr.InsertIdErr):
			helpers.ClientErr(w, http.StatusBadRequest)
			// errors from POST patient methods
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

func logIn(w http.ResponseWriter, r *http.Request) error {
	var login dto.LogIn

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&login); err != nil {
		return err
	}

	if err := validateLoginStruct(login); err != nil {
		return err
	}

	role := *login.RoleId
	loginStr := strings.TrimSpace(*login.Login)
	password := strings.TrimSpace(*login.Password)

	*login.Patient.Login = loginStr

	err := validateLogin(loginStr, role)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = validatePassword(password)

	if err != nil {
		return nil
	}

	//password, err = hashPassword(password)

	*login.Patient.Password = password

	if err != nil {
		return err
	}

	switch {
	case role == 1 || role == 3:
		if login.Patient == nil {
			return NoPatientStructErr
		}

		if err := logInPatientAdmin(w, r, login); err != nil {
			return err
		}
	case role == 2:
		if login.Doctor == nil {
			return NoDoctorStructErr
		}

		if err := logInDoctor(w, role, loginStr, password, *login.Doctor); err != nil {
			return err
		}
	default:
		return NoRoleErr
	}

	return nil
}

func logInPatientAdmin(w http.ResponseWriter, r *http.Request, login dto.LogIn) error {
	role := *login.RoleId
	patient := *login.Patient

	loginStr := *patient.Login
	password := *patient.Password

	body, err := json.Marshal(patient)

	if err != nil {
		return err
	}

	var b bytes.Buffer
	b.Write(body)
	r.Body = io.NopCloser(&b)

	err = patients.PostPatientReturnErr(w, r)

	if err != nil {
		return err
	}

	id, err := postgres.DataBase.CheckLoginAndPasswordPatient(loginStr, password)

	if err != nil {
		return err
	}

	accessToken, refreshToken, err := jwtgen.GeneratePairsToken(role, id)

	if err != nil {
		return err
	}

	err = postgres.DataBase.UpdatePatientRefreshToken(id, refreshToken)

	if err != nil {
		return err
	}

	jwtokens := dto.NewJWTokens(accessToken, refreshToken)
	tokens, err := jwtokens.EncodeToJSON()

	if err != nil {
		return err
	}

	w.Write(tokens)
	return nil
}

func logInDoctor(w http.ResponseWriter, role int, loginStr, password string, doctor postgres.Doctor) error {
	*doctor.Login = loginStr
	*doctor.Password = password
	err := postDoctorFromStruct(doctor)

	if err != nil {
		return err
	}

	id, err := postgres.DataBase.CheckLoginAndPasswordDoctor(loginStr, password)

	if err != nil {
		return err
	}

	accessToken, refreshToken, err := jwtgen.GeneratePairsToken(role, id)

	if err != nil {
		return err
	}

	err = postgres.DataBase.UpdatePatientRefreshToken(id, refreshToken)

	if err != nil {
		return err
	}

	jwtokens := dto.NewJWTokens(accessToken, refreshToken)
	tokens, err := jwtokens.EncodeToJSON()

	if err != nil {
		return err
	}

	w.Write(tokens)
	return nil
}

func postPatientFromStruct(patient postgres.Patient) error {
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

	return nil
}

func postDoctorFromStruct(doctor postgres.Doctor) error {
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

	return nil
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	err := signIn(w, r)

	if err != nil {
		logger.ErrLog.Println(err)
		switch {
		case errors.Is(err, NoRoleErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, NoLoginErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, NoPasswordErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		case errors.Is(err, IncorrectLoginOrPasswordErr):
			helpers.ClientErr(w, http.StatusBadRequest)
		default:
			helpers.ServeErr(w, err)

		}
	}
}

func signIn(w http.ResponseWriter, r *http.Request) error {
	var login dto.LogIn

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&login); err != nil {
		return err
	}

	if err := validateLoginStruct(login); err != nil {
		return err
	}

	role := *login.RoleId
	loginStr := strings.TrimSpace(*login.Login)
	password := strings.TrimSpace(*login.Password)

	switch {
	case role == 1 || role == 3:
		if err := signInPatientAdmin(w, role, loginStr, password); err != nil {
			return err
		}
	case role == 2:
		if err := signInDoctor(w, role, loginStr, password); err != nil {
			return err
		}
	default:
		return NoRoleErr
	}

	return nil
}

func signInPatientAdmin(w http.ResponseWriter, role int, loginStr, password string) error {
	id, err := postgres.DataBase.CheckLoginAndPasswordPatient(loginStr, password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return IncorrectLoginOrPasswordErr
		}

		return err
	}

	accessToken, refreshToken, err := jwtgen.GeneratePairsToken(role, id)

	if err != nil {
		return err
	}

	err = postgres.DataBase.UpdatePatientRefreshToken(id, refreshToken)

	if err != nil {
		return err
	}

	jwtokens := dto.NewJWTokens(accessToken, refreshToken)
	tokens, err := jwtokens.EncodeToJSON()

	if err != nil {
		return err
	}

	w.Write(tokens)
	return nil
}

func signInDoctor(w http.ResponseWriter, role int, loginStr, password string) error {
	id, err := postgres.DataBase.CheckLoginAndPasswordDoctor(loginStr, password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return IncorrectLoginOrPasswordErr
		}

		return err
	}

	accessToken, refreshToken, err := jwtgen.GeneratePairsToken(role, id)

	if err != nil {
		return err
	}

	err = postgres.DataBase.UpdateDoctorRefreshToken(id, refreshToken)

	if err != nil {
		return err
	}

	jwtokens := dto.NewJWTokens(accessToken, refreshToken)
	tokens, err := jwtokens.EncodeToJSON()

	if err != nil {
		return err
	}

	w.Write(tokens)

	return nil
}

func validateLoginStruct(login dto.LogIn) error {
	if login.RoleId == nil {
		return NoRoleErr
	}

	if login.Login == nil {
		return NoLoginErr
	}

	if login.Password == nil {
		return NoPasswordErr
	}

	return nil
}

func validateLogin(login string, role int) error {
	login, err := validateLoginPasswordLength(login)

	if err != nil {
		return err
	}

	switch {
	case role == 1 || role == 3:
		err = postgres.DataBase.CheckLoginPatient(login)
	case role == 2:
		err = postgres.DataBase.CheckLoginDoctor(login)
	default:
		return NoRoleErr
	}

	return err
}

func validatePassword(password string) error {
	password, err := validateLoginPasswordLength(password)

	if err != nil {
		return err
	}

	password, err = hashPassword(password)

	if err != nil {
		return err
	}

	return nil
}

func validateLoginPasswordLength(word string) (string, error) {
	word = strings.TrimSpace(word)

	if utf8.RuneCountInString(word) < 8 {
		return "", LenLoginPasswordErr
	}

	return word, nil
}

func hashPassword(password string) (string, error) {
	return hash(password)
}

func hash(word string) (string, error) {
	hash := sha1.New()

	if _, err := io.WriteString(hash, word); err != nil {
		return "", err
	}

	return string(hash.Sum(nil)), nil
}
