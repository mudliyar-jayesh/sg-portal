package util

import (
	"context"
)

type key int

const userKey key = 0

// ContextWithUserID stores the user ID in the context.
func ContextWithUserID(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, userKey, userID)
}

// UserIDFromContext retrieves the user ID from the context.
func UserIDFromContext(ctx context.Context) (uint64, bool) {
	userID, ok := ctx.Value(userKey).(uint64)
	return userID, ok
}

