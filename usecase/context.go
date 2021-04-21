package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const userIDCtxKey = "userID"

func (u *Usecase) UUIDToCtx(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDCtxKey, userID)
}

func UUIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(userIDCtxKey)
	valUUID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("failed to fetch id from ctx")
	}
	return valUUID, nil
}
