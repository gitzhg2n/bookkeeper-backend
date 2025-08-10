package routes

import (
	"context"

	"bookkeeper-backend/middleware"
)

func mockUserContext(ctx context.Context) context.Context {
	return middleware.WithUser(ctx, &middleware.UserContext{
		ID:    1,
		Email: "test@example.com",
		Role:  "user",
	})
}

func contextWithUserID(ctx context.Context, userID uint) context.Context {
	return middleware.WithUser(ctx, &middleware.UserContext{
		ID:    userID,
		Email: "test@example.com",
		Role:  "user",
	})
}

func contextWithRole(ctx context.Context, role string) context.Context {
	return middleware.WithUser(ctx, &middleware.UserContext{
		ID:    1,
		Email: "test@example.com",
		Role:  role,
	})
}