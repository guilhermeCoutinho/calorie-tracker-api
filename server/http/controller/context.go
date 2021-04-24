package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/sirupsen/logrus"
)

const ctxKey = "ctxKey"

func ClaimsToCtx(ctx context.Context, claims *Claims) context.Context {
	rawBytes, _ := json.Marshal(claims)
	return context.WithValue(ctx, ctxKey, rawBytes)
}

func ClaimsFromCtx(ctx context.Context) (*Claims, error) {
	val, ok := ctx.Value(ctxKey).([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to assert context")
	}
	var claim Claims
	err := json.Unmarshal(val, &claim)
	if err != nil {
		return nil, err
	}

	return &claim, nil
}

func LoggerFromCtx(ctx context.Context) logrus.FieldLogger {
	return ctx.Value(models.LoggerCtxKey).(logrus.FieldLogger)
}
