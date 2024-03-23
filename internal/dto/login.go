package dto

import (
	"encoding/json"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"net/http"
)

type LogIn struct {
	RoleId   *int              `json:"role_id"`
	Login    *string           `json:"login"`
	Password *string           `json:"password"`
	Patient  *postgres.Patient `json:"patient"`
	Doctor   *postgres.Doctor  `json:"doctor"`
}

func (l *LogIn) DecodeFromJSON(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(l); err != nil {
		return err
	}

	return nil
}
