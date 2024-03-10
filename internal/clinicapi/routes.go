package clinicapi

import (
	"github.com/Vadim992/clinicAPI/internal/clinicapi/requests/appointments"
	"github.com/Vadim992/clinicAPI/internal/clinicapi/requests/doctors"
	"github.com/Vadim992/clinicAPI/internal/clinicapi/requests/patients"
	"github.com/Vadim992/clinicAPI/internal/helpers"
	"github.com/Vadim992/clinicAPI/internal/middlewares"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func Routes() http.Handler {

	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(helpers.RouterNotFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(helpers.RouterMethodNotAllowed)

	setPatientControllers(router)
	setDoctorControllers(router)
	setAppointmentControllers(router)

	mwChain := alice.New(middlewares.RecoverPanic, middlewares.LogRequest, middlewares.SetHeaders)

	return mwChain.Then(router)
}

func setPatientControllers(router *mux.Router) {
	router.HandleFunc("/patients", patients.GetPatients).Methods("GET")
	router.HandleFunc("/patients", patients.PostPatient).Methods("POST")

	router.HandleFunc("/patients/{id:[0-9]+}", patients.GetPatientId).Methods("GET")
	router.HandleFunc("/patients/{id:[0-9]+}", patients.PatchPatientId).Methods("PATCH")
	router.HandleFunc("/patients/{id:[0-9]+}", patients.DeletePatientId).Methods("DELETE")
}

func setDoctorControllers(router *mux.Router) {
	router.HandleFunc("/doctors", doctors.GetDoctors).Methods("GET")
	router.HandleFunc("/doctors", doctors.PostDoctor).Methods("POST")

	router.HandleFunc("/doctors/{id:[0-9]+}", doctors.GetDoctorId).Methods("GET")
	router.HandleFunc("/doctors/{id:[0-9]+}", doctors.PatchDoctorId).Methods("PATCH")
	router.HandleFunc("/doctors/{id:[0-9]+}", doctors.DeleteDoctorId).Methods("DELETE")
}

func setAppointmentControllers(router *mux.Router) {
	router.HandleFunc("/appointments", appointments.GetAppointments).Methods("GET")
	router.HandleFunc("/appointments", appointments.PostAppointment).Methods("POST")

	router.HandleFunc("/appointments/{id:[0-9]+}", appointments.GetAppointmentId).Methods("GET")
	router.HandleFunc("/appointments/{id:[0-9]+}", appointments.PatchAppointmentId).Methods("PATCH")
	router.HandleFunc("/appointments/{id:[0-9]+}", appointments.DeleteAppointmentId).Methods("DELETE")

	router.HandleFunc("/appointments/per_doctor",
		appointments.GetAppointmentsCountAllDoctors).Methods("GET")
}
