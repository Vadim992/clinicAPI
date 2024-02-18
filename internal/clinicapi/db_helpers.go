package clinicapi

import (
	"database/sql"
	"fmt"
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"github.com/Vadim992/clinicAPI/pkg/validator/validate"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//checkEmail checks Is valid PATIENT's email or not

func checkEmail(email sql.NullString) bool {
	if email.Valid {
		return validate.ValidateEmail(email.String)
	}

	return true
}

// checkStruct checks Is valid struct (Patient or Doctor) or not
func checkStructs(structs interface{}) bool {
	switch structs.(type) {
	case postgres.Patient:
	case postgres.Doctor:
	case postgres.Record:

	default:
		return false
	}

	v := reflect.ValueOf(structs)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Name == "Id" {
			continue
		}

		switch v.Field(i).Type().Kind() {
		case reflect.Int:
			if v.Field(i).Int() == 0 {
				return false
			}
		case reflect.String:

			if v.Field(i).String() == "" {
				return false
			}
		case reflect.Struct:
			structVal := v.Field(i).Interface()

			switch curStruct := structVal.(type) {
			case time.Time:
				if curStruct.IsZero() {
					return false
				}
			}

		}
	}

	return true
}

func patchStructs(structs interface{}) string {
	switch structs.(type) {
	case postgres.Patient:
	case postgres.Doctor:
	default:
		return ""
	}

	v := reflect.ValueOf(structs)
	t := v.Type()

	var res strings.Builder

	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Name == "Id" {
			continue
		}

		switch v.Field(i).Type().Kind() {
		case reflect.Int:
			nameField := t.Field(i).Name
			val := v.Field(i).Int()

			if val != 0 {
				res.WriteString(fmt.Sprintf("%s=%d, ", nameField, val))
			}

		case reflect.String:

			nameField := t.Field(i).Name
			val := v.Field(i).String()

			if val != "" {
				res.WriteString(fmt.Sprintf("%s='%s', ", nameField, val))
			}

		case reflect.Struct:
			nameField := t.Field(i).Name
			structVal := v.Field(i).Interface()

			switch curStruct := structVal.(type) {
			case sql.NullString:

				if curStruct.String != "" {
					res.WriteString(fmt.Sprintf("%s='%s', ", nameField, curStruct.String))
				}

			}

		}
	}

	str := strings.TrimSpace(res.String())
	str = str[:len(str)-1]

	return str
}

// func checkPhoneNum checks Is valid PATIENT's phone number or not
func checkPhoneNum(phoneNum string) bool {
	return validate.ValidatePhoneNum(phoneNum)
}

// formatPhoneNumber formats every phone number to 81112223344
func formatPhoneNumber(phoneNum string) string {
	var formatedPhone strings.Builder

	for i := 0; i < len(phoneNum); i++ {
		if phoneNum[i] == '+' {
			formatedPhone.WriteByte('8')
			i++
			continue
		}

		if _, err := strconv.Atoi(string(phoneNum[i])); err == nil {
			formatedPhone.WriteByte(phoneNum[i])
		}
	}
	return formatedPhone.String()
}

// validateNumPage checks that page > 0
func validateNumPage(page int) bool {
	if page < 1 {
		return false
	}
	return true
}
