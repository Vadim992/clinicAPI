package clinicapi

import "net/http"

type Decoder interface {
	Decode(r *http.Request) error
}
