package middleware

import "context"

type userContextKeyType struct{}

var userContextKey = userContextKeyType{}

type UserContext struct {
	ID    uint
	Email string
	Role  string
	// Future: household membership, plan, etc.
	Plan  string // free, premium, selfhost
}

func WithUser(ctx context.Context, u *UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

func UserFrom(ctx context.Context) (*UserContext, bool) {
	u, ok := ctx.Value(userContextKey).(*UserContext)
	return u, ok
}