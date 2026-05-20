package reqctx

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/google/uuid"
)

type contextKey string

const callerIDKey contextKey = "callerID"

func WithCallerID(ctx context.Context, callerID uuid.UUID) context.Context {
	return context.WithValue(ctx, callerIDKey, callerID)
}

func CallerIDFrom(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(callerIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, auth.ErrUnauthenticated
	}

	return id, nil
}
