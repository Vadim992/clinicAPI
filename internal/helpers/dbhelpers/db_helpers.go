package dbhelpers

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func ConvertIdFromStrToInt(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// FormatPhoneNumber formats every phone number to 81112223344
func FormatPhoneNumber(phoneNum string) string {
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
