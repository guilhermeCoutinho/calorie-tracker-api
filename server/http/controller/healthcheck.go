package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/messages"
)

type HealthCheck struct{}

func NewHealthcheck() *HealthCheck {
	return &HealthCheck{}
}

func (m *HealthCheck) GetPing(ctx context.Context, args *struct{}) (*messages.BaseResponse, error) {
	return &messages.BaseResponse{Msg: "Pong", Code: http.StatusOK}, fmt.Errorf("Testing error new")
}
