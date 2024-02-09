package clinicapi

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

func (c *ClinicAPI) serveErr(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%e\n%s", err, debug.Stack())

	if err := c.ErrLog.Output(2, trace); err != nil {
		c.ErrLog.Println("failed to show stack trace")
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (c *ClinicAPI) clientErr(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (c *ClinicAPI) notFound(w http.ResponseWriter) {
	c.clientErr(w, http.StatusNotFound)
}

// Add functions for routers
func (c *ClinicAPI) routerNotFound(w http.ResponseWriter, r *http.Request) {
	c.notFound(w)
}

func (c *ClinicAPI) routerMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	c.clientErr(w, http.StatusMethodNotAllowed)
}

// Get env var from .env file

func GetProjectEnv() (host_db, port_db, username_db, password_db, dbname_db, searchPath_db string) {
	host_db = os.Getenv("HOST_DB")
	port_db = os.Getenv("PORT_DB")
	username_db = os.Getenv("USERNAME_DB")
	password_db = os.Getenv("PASSWORD_DB")
	dbname_db = os.Getenv("DBNAME_DB")
	searchPath_db = os.Getenv("SEARCHPATH_DB")

	return
}
