package clinicapi

import (
	"net/http"
	"strconv"
	"strings"
)

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

func decode(r *http.Request, decoder Decoder) error {
	return decoder.Decode(r)
}
