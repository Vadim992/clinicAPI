package helpers

import (
	"fmt"
	"github.com/Vadim992/clinicAPI/pkg/logger"
	"net/http"
	"runtime/debug"
)

func ServeErr(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%e\n%s", err, debug.Stack())

	if err := logger.ErrLog.Output(2, trace); err != nil {
		logger.ErrLog.Println("failed to show stack trace")
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func ClientErr(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func NotFound(w http.ResponseWriter) {
	ClientErr(w, http.StatusNotFound)
}

// routerNotFound adds functions for routers
func RouterNotFound(w http.ResponseWriter, r *http.Request) {
	NotFound(w)
}

func RouterMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	ClientErr(w, http.StatusMethodNotAllowed)
}
