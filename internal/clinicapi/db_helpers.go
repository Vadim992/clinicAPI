package clinicapi

import (
	"database/sql"
	"fmt"
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"github.com/Vadim992/clinicAPI/pkg/validator/validate"
	"reflect"
	"strings"
)

//checkEmail checks Is valid PATIENT email or not

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
