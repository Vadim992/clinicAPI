package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (c *clinicAPI) serveErr(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%e\n%s", err, debug.Stack())

	if err := c.errLog.Output(2, trace); err != nil {
		c.errLog.Println("failed to show stack trace")
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (c *clinicAPI) clientErr(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (c *clinicAPI) notFound(w http.ResponseWriter) {
	c.clientErr(w, http.StatusNotFound)
}
