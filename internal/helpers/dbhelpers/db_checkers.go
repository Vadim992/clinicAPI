package dbhelpers

import (
	"fmt"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"github.com/Vadim992/clinicAPI/internal/helpers/customerr/structserr"
	"reflect"
	"strings"
	"time"
)

// CheckStructsFields checks Is valid struct (Patient or Doctor) or not
func CheckStructsFields(structs interface{}) error {
	canNilFields, isValidStruct := validateStruct(structs)

	if !isValidStruct {
		return structserr.InvalidTypeOfStructErr
	}

	v := reflect.ValueOf(structs)
	t := v.Type()

	isSkipId := skipId(structs)

	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Name == "Id" && isSkipId {
			continue
		}

		if v.Field(i).IsNil() {
			if _, ok := canNilFields[t.Field(i).Name]; !ok {
				return structserr.EmptyFieldErr
			}
		}
	}

	return nil
}

func validateStruct(structs interface{}) (map[string]struct{}, bool) {
	m := make(map[string]struct{})

	switch structs.(type) {
	case postgres.Patient:
		m["Email"] = struct{}{}
		m["RefreshToken"] = struct{}{}
	case postgres.Doctor:
		m["RefreshToken"] = struct{}{}
	case postgres.Record:
		m["PatientId"] = struct{}{}
	default:
		return m, false
	}

	return m, true
}

func skipId(structs interface{}) (isSkip bool) {
	switch structs.(type) {
	case postgres.Patient:
		isSkip = true
	case postgres.Doctor:
		isSkip = true
	}

	return
}

// checkTypeStructs checks struct in CASE reflect.Struct of  putCheckStructs func
//func putCheckTypeStructs(structVal interface{}) bool {
//	switch curStruct := structVal.(type) {
//	case time.Time:
//		if curStruct.IsZero() {
//			return false
//		}
//	}
//
//	return true
//}

// CheckStructFieldsPatch return request string for PATCH, IMPORTANT: before CheckStructFieldsPatch you MUST VALIDATE
// your struct
func CheckStructFieldsPatch(structs interface{}) (string, error) {
	_, isValidStruct := validateStruct(structs)
	if !isValidStruct {
		return "", structserr.InvalidTypeOfStructErr
	}

	v := reflect.ValueOf(structs)
	t := v.Type()

	isSkipId := skipId(structs)

	var res strings.Builder

	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Name == "Id" && isSkipId {
			continue
		}

		if !v.Field(i).IsNil() {
			fieldType := v.Field(i).Interface()

			switch currentType := fieldType.(type) {
			case *string:
				nameField := t.Field(i).Name
				val := *currentType

				res.WriteString(fmt.Sprintf("%s='%s', ", nameField, val))
			case *int:
				nameField := t.Field(i).Name
				val := *currentType

				res.WriteString(fmt.Sprintf("%s=%d, ", nameField, val))
			case *time.Time:
				nameField := t.Field(i).Name
				val := *currentType

				res.WriteString(fmt.Sprintf("%s='%s', ", nameField, val.Format(time.DateTime)))
			}
		}
	}

	str := strings.TrimSpace(res.String())

	if len(str) != 0 {
		str = str[:len(str)-1]
	}

	return str, nil
}

//func patchCheckTypeStructs(structVal interface{}, nameField string, b *strings.Builder) {
//
//	switch curStruct := structVal.(type) {
//	case sql.NullString:
//
//		if curStruct.Valid {
//			b.WriteString(fmt.Sprintf("%s='%s', ", nameField, curStruct.String))
//		}
//	case sql.NullInt64:
//
//		if curStruct.Valid {
//			b.WriteString(fmt.Sprintf("%s='%d', ", nameField, curStruct.Int64))
//		}
//
//	case time.Time:
//
//		if !curStruct.IsZero() {
//			b.WriteString(fmt.Sprintf("%s='%v', ", nameField, curStruct.Format("2006-01-02 15:04:05")))
//		}
//
//	}
//
//}
