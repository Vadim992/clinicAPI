package clinicapi

import (
	"fmt"
	"net/http"
)

func setHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func (c *ClinicAPI) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.InfoLog.Printf("%s %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})

}

func (c *ClinicAPI) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				w.Header().Set("Connection", "close")

				c.serveErr(w, fmt.Errorf("%s", err))
			}

			next.ServeHTTP(w, r)
		}()
	})

}
