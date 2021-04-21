package usecase

import (
	"context"
	"encoding/json"
	"fmt"
)

const ctxKey = "ctxKey"

func (u *Usecase) ClaimsToCtx(ctx context.Context, claims *Claims) context.Context {
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
