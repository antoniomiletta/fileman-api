package reqctx

import (
	"context"
	"fmt"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/google/uuid"
)

type ctxKey int

const callerIDKey ctxKey = iota

func WithCallerID(ctx context.Context, callerID uuid.UUID) context.Context {
	return context.WithValue(ctx, callerIDKey, callerID)
}

func CallerIDFrom(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(callerIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("%w: missing token", auth.ErrUnauthenticated)
	}

	return id, nil
}
