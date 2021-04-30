package healthcheck

import (
	"context"

	"github.com/guilhermeCoutinho/api-studies/controller/contextextensions"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
)

type HealthCheck struct{}

func NewHealthcheck() *HealthCheck {
	return &HealthCheck{}
}

func (m *HealthCheck) Get(ctx context.Context, args *struct{}, vars *struct{}) (*struct{}, *wrapper.HandlerError) {
	logger := contextextensions.LoggerFromCtx(ctx)
	logger.Info("Pong")
	return &struct{}{}, nil
}
