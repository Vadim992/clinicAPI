package clinicapi

import (
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"log"
)

type ClinicAPI struct {
	ErrLog  *log.Logger
	InfoLog *log.Logger
	DB      *postgres.DB
}
