package reqctx

import (
	"context"
	"fmt"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/google/uuid"
)

type ctxKey int

const callerIDKey ctxKey = iota

// WithCallerID returns a derived context with the provided callerID value.
func WithCallerID(ctx context.Context, callerID uuid.UUID) context.Context {
	return context.WithValue(ctx, callerIDKey, callerID)
}

// CallerIDFrom returns the callerID value from the provided context.
// If no value is found, it returns a wrapped auth.ErrUnauthenticated.
func CallerIDFrom(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(callerIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("missing token: %w", auth.ErrUnauthenticated)
	}

	return id, nil
}
