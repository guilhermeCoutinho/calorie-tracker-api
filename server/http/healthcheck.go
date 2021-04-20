package http

import (
	"net/http"
)

type healthcheckController struct{}

func newHealthcheckController() *healthcheckController {
	return &healthcheckController{}
}

func (m *healthcheckController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WORKING"))
}
