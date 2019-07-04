package http

import "net/http"

type Applier interface {
	ApplyAuth(next http.Handler) http.Handler
	PreFlightHandler(w http.ResponseWriter, r *http.Request)
}
