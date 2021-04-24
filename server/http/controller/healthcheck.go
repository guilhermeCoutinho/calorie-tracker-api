package controller

import (
	"context"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/messages"
)

type HealthCheck struct{}

func NewHealthcheck() *HealthCheck {
	return &HealthCheck{}
}

func (m *HealthCheck) Get(ctx context.Context, args *struct{}) (*messages.BaseResponse, error) {
	logger := LoggerFromCtx(ctx)
	logger.Info("Pong")
	return &messages.BaseResponse{Msg: "Pong", Code: http.StatusOK}, nil
}
