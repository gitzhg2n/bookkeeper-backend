# Migration Notes - Stage 1 Backend Consolidation

## Overview
This document describes the Stage 1 consolidation of the bookkeeper-backend to a single `/v1` API stack, removing legacy code and implementing security improvements.

## Changes Made

### 1. Legacy Code Removal
- **Legacy models package**: No legacy `models/` package was found to delete - models are properly organized in `internal/models/`
- **Legacy routes**: No legacy route handlers were found to delete - routes are properly organized with `/v1` prefix
- **Legacy middleware**: No legacy middleware files were found to delete - middleware is properly organized
- **Legacy entrypoint**: No root `main.go` was found - application properly uses `cmd/server/main.go`

### 2. Dependency Cleanup
- **gorilla/mux removal**: Removed unused gorilla/mux dependency from go.mod and go.sum via `go mod tidy`
- All routing is handled by the standard `net/http` ServeMux

### 3. Security Improvements

#### KEK Derivation
- **New file**: `internal/security/kek.go`
- Added `DeriveKEK()` function using HMAC-SHA256 expansion for proper key separation
- Updated `routes/auth.go` Register function to derive KEK using `security.DeriveKEK(passwordKey, "bookkeeper:dek:v1")`
- Password hash (Argon2 output) is still stored separately from KEK for verification

#### Rate Limiting
- Added rate limiting middleware to auth endpoints: `/v1/auth/register`, `/v1/auth/login`, `/v1/auth/refresh`, `/v1/auth/logout`
- Rate limit: 10 requests per 60 seconds per IP address
- Rate limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`) are included in responses

### 4. API Enhancements
- **New endpoint**: `DELETE /v1/households/{id}/budgets/{budgetID}` 
- Route properly wired to `BudgetHandler.Delete()` method
- Returns 200 status and removes the budget record

### 5. Code Organization
- **Response helpers**: Already centralized in `routes/response_helpers.go` - no duplicates found
- **Validation helpers**: Already centralized in `routes/validation.go` (sanitizeString, etc.)
- **Shared utilities**: Properly organized with no duplication

### 6. Testing
- **New integration test**: `internal/tests/core_flow_test.go`
- Tests complete vertical slice: register → create household → create category → create account → upsert budget → create transaction → list budgets
- Verifies `actual_cents` reflects created transactions in budget listings

### 7. Infrastructure
- **Dockerfile updates**: 
  - Build target changed from `./main.go` to `./cmd/server`
  - Removed unnecessary `.env` file copying
- **Migration placeholder**: Ready for future `0003_future_security_columns.sql` if needed

## Files Modified
- `routes/auth.go` - Added KEK derivation using new security helper
- `routes/router.go` - Added rate limiting to auth endpoints and DELETE budget route
- `Dockerfile` - Updated build target and removed .env copying
- `go.mod`/`go.sum` - Cleaned up dependencies (removed gorilla/mux)

## Files Added
- `internal/security/kek.go` - KEK derivation helper functions
- `internal/tests/core_flow_test.go` - End-to-end integration test
- `MIGRATION_NOTES.md` - This documentation

## Current State
The application is now consolidated to a single `/v1` API stack with:
- ✅ Proper key derivation separation (KEK vs password verification)
- ✅ Rate limiting on authentication endpoints
- ✅ Complete CRUD operations for budgets including DELETE
- ✅ Comprehensive integration testing
- ✅ Clean dependency management
- ✅ Modern containerization setup

## Next Steps (Future Stages)
1. **Password verifier separation**: Move from storing Argon2 output as password hash to a separate password verifier column
2. **Refresh token reuse detection**: Implement token family rotation with reuse detection
3. **Enhanced encryption**: Full data-at-rest encryption using derived DEKs
4. **Audit logging**: Comprehensive audit trail for sensitive operations
5. **Multi-factor authentication**: TOTP/WebAuthn support

## Verification Commands
```bash
# Verify no legacy model imports
grep -R "\"bookkeeper-backend/models\"" . || echo "✅ No legacy model imports"

# Verify no gorilla/mux usage  
grep -R 'gorilla/mux' . || echo "✅ No gorilla/mux dependency"

# Build verification
go build ./... && echo "✅ Build successful"

# Test verification
go test ./... && echo "✅ All tests pass"

# Docker build verification
docker build -t bookkeeper-backend . && echo "✅ Docker build successful"
```

## Database Schema
Current schema supports the consolidated structure with proper foreign key relationships and no breaking changes were required.

**Note**: This Stage 1 consolidation found the codebase was already well-organized with proper `/v1` structure. The main improvements were in security (KEK derivation, rate limiting), testing (integration test), and cleanup (dependency management).