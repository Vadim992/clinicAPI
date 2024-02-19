package clinicapi

import (
	"database/sql"
	"fmt"
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"reflect"
	"strings"
	"time"
)

// checkStruct checks Is valid struct (Patient or Doctor) or not
func putCheckStructs(structs interface{}) bool {

	if !validateStruct(structs) {
		return false
	}

	v := reflect.ValueOf(structs)
	t := v.Type()

	isSkipId := skipId(structs)

	for i := 0; i < v.NumField(); i++ {

		if t.Field(i).Name == "Id" && isSkipId {
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

			if !putCheckTypeStructs(structVal) {
				return false
			}
		}
	}

	return true
}

func validateStruct(structs interface{}) bool {
	switch structs.(type) {
	case postgres.Patient:
	case postgres.Doctor:
	case postgres.Record:

	default:
		return false
	}
	return true
}

func skipId(structs interface{}) (isSkip bool) {
	isSkip = false
	switch structs.(type) {
	case postgres.Patient:
		isSkip = true
	case postgres.Doctor:
		isSkip = true
	}
	return
}

// checkTypeStructs checks struct in CASE reflect.Struct of  putCheckStructs func
func putCheckTypeStructs(structVal interface{}) bool {

	switch curStruct := structVal.(type) {

	case time.Time:
		if curStruct.IsZero() {
			return false
		}
	}

	return true
}

// patchCheckStructs return request string for PATCH, IMPORTANT: before patchCheckStructs you MUST VALIDATE
// your struct
func patchCheckStructs(structs interface{}) string {

	if !validateStruct(structs) {
		return ""
	}

	v := reflect.ValueOf(structs)
	t := v.Type()

	isSkipId := skipId(structs)

	var res strings.Builder

	for i := 0; i < v.NumField(); i++ {

		if t.Field(i).Name == "Id" && isSkipId {
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

			patchCheckTypeStructs(structVal, nameField, &res)

		}
	}

	str := strings.TrimSpace(res.String())

	if len(str) != 0 {
		str = str[:len(str)-1]
	}

	return str
}

func patchCheckTypeStructs(structVal interface{}, nameField string, b *strings.Builder) {

	switch curStruct := structVal.(type) {
	case sql.NullString:

		if curStruct.Valid {
			b.WriteString(fmt.Sprintf("%s='%s', ", nameField, curStruct.String))
		}
	case sql.NullInt64:

		if curStruct.Valid {
			b.WriteString(fmt.Sprintf("%s='%d', ", nameField, curStruct.Int64))
		}

	case time.Time:

		if !curStruct.IsZero() {
			b.WriteString(fmt.Sprintf("%s='%v', ", nameField, curStruct.Format("2006-01-02 15:04:05")))
		}

	}

}
