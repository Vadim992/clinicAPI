package clinicapi

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func (c *ClinicAPI) Routes() http.Handler {

	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(c.routerNotFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(c.routerMethodNotAllowed)

	router.HandleFunc("/patients", c.handlePatients).Methods("GET", "POST")
	router.HandleFunc("/patients/{id:[0-9]+}", c.handlePatientId).Methods("GET", "PUT", "PATCH", "DELETE")

	router.HandleFunc("/doctors", c.handleDoctors).Methods("GET", "POST")
	router.HandleFunc("/doctors/{id:[0-9]+}", c.handleDoctorId).Methods("GET", "PUT", "PATCH", "DELETE")

	mwChain := alice.New(c.recoverPanic, c.logRequest, setHeaders)

	return mwChain.Then(router)
}
