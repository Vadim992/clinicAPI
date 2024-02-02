package main

import "net/http"

func (c *clinicAPI) routes() http.Handler {

	router := http.NewServeMux()

	router.HandleFunc("/client/", c.handleClient)
	router.HandleFunc("/doctor/", c.handleDoctor)

	return c.recoverPanic(c.logRequest(setHeaders(router)))
}
