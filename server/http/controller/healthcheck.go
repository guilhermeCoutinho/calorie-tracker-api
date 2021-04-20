package controller

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type HealthCheck struct {
	logger logrus.FieldLogger
}

func NewHealthcheck(
	logger logrus.FieldLogger,
) *HealthCheck {
	return &HealthCheck{
		logger: logger,
	}
}

func (m *HealthCheck) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WORKING"))
}
