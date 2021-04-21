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
	m.logger.Debug("Working")
	w.Write([]byte("Pong"))
}
