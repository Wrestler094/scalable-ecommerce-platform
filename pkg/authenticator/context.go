package authenticator

import (
	"context"
)

type userIDKey struct{}
type roleKey struct{}

func WithUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

func UserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey{}).(int64)
	return id, ok
}

func WithUserRole(ctx context.Context, role Role) context.Context {
	return context.WithValue(ctx, roleKey{}, role)
}

func UserRole(ctx context.Context) (Role, bool) {
	role, ok := ctx.Value(roleKey{}).(Role)
	return role, ok
}
