# Middleware Utilities (Go Version)

This folder contains reusable middleware for privacy, security, and resource access in Bookkeeper.

## Ownership Middleware
- Ensures users can only access resources they own (households, accounts, budgets, goals, investments, transactions).
- See: `ownership.go`

## Rate Limit Middleware
- Throttles requests to prevent brute-force attacks.
- Upgrade to Redis/memcached for production.
- See: `rate_limit.go`

## Role Middleware
- Restricts sensitive actions (user create/delete, etc.) to admins only.
- See: `role.go`

## Adding Middleware
All middleware is designed to be composable and used in Go HTTP router handlers.

```go
import "bookkeeper-backend-go/middleware"

if !middleware.CheckAccountOwnership(ctx, userID, accountID) {
    // return forbidden
}
```