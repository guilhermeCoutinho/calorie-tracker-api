package controller

import (
	"context"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
)

type HealthCheck struct{}

func NewHealthcheck() *HealthCheck {
	return &HealthCheck{}
}

func (m *HealthCheck) Get(ctx context.Context, args *struct{}, vars *struct{}) (*messages.BaseResponse, *wrapper.HandlerError) {
	logger := LoggerFromCtx(ctx)
	logger.Info("Pong")
	return &messages.BaseResponse{Msg: "Pong", Code: http.StatusOK}, nil
}
