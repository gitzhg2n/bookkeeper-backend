package routes

import (
	"context"

	"bookkeeper-backend/middleware"
)

// mockUserContext creates a mock user context for testing
func mockUserContext(ctx context.Context) context.Context {
	userCtx := &middleware.UserContext{
		ID:           1,
		Email:        "test@example.com",
		Role:         "user",
		HouseholdIDs: []uint{1},
		AccountIDs:   []uint{1},
	}
	return context.WithValue(ctx, "userContext", userCtx)
}

// contextWithUserID creates a context with a specific user ID for testing
func contextWithUserID(ctx context.Context, userID uint) context.Context {
	userCtx := &middleware.UserContext{
		ID:           userID,
		Email:        "test@example.com",
		Role:         "user",
		HouseholdIDs: []uint{1},
		AccountIDs:   []uint{1},
	}
	return context.WithValue(ctx, "userContext", userCtx)
}

// contextWithRole creates a context with a specific role for testing
func contextWithRole(ctx context.Context, role string) context.Context {
	userCtx := &middleware.UserContext{
		ID:           1,
		Email:        "test@example.com",
		Role:         role,
		HouseholdIDs: []uint{1},
		AccountIDs:   []uint{1},
	}
	return context.WithValue(ctx, "userContext", userCtx)
}