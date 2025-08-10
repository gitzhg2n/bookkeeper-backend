Migration Notes - Stage 1 Backend Consolidation

Overview
This document summarizes the Stage 1 consolidation of the bookkeeper-backend to a single /v1 API stack, removing legacy code, cleaning dependencies, and introducing security improvements.

Changes Made

Legacy Code and Structure

Legacy models: No legacy models/ package was found to delete; models are correctly in internal/models/.

Legacy routes: No legacy route handlers were found to delete; routes are organized under the /v1 prefix.

Legacy middleware: No legacy middleware files were found to delete; middleware is properly organized.

Legacy entrypoint: No root main.go was found; the application uses cmd/server/main.go.

Dependency Cleanup

Removed unused gorilla/mux dependency via go mod tidy.

All routing is handled by net/http ServeMux.

Security Improvements

KEK Derivation

New file: internal/security/kek.go

Added DeriveKEK() using HMAC-SHA256 expansion for proper key separation.

routes/auth.go Register updated to use security.DeriveKEK(passwordKey, "bookkeeper:dek:v1").

KEK is derived separately from the password hash; future work will further separate password verification data.

Rate Limiting

Added rate limiting to auth endpoints: /v1/auth/register, /v1/auth/login, /v1/auth/refresh, /v1/auth/logout.

Limit: 10 requests per 60s per IP.

Responses include X-RateLimit-Limit and X-RateLimit-Remaining headers.

API Enhancements

New endpoint: DELETE /v1/households/{id}/budgets/{budgetID}

Route wired in routes/router.go; handler deletes the budget and returns success.

Code Organization and Cleanup

Centralized small shared utilities in routes/util.go (e.g., parseUintString) to reduce duplication.

Cleaned imports, fixed incorrect paths, and removed unused imports.

Dockerfile updated to build from ./cmd/server and removed .env copy.

Placeholder migration added: internal/db/migrations/0003_future_security_columns.sql for future password_verifier column (non-breaking).

Testing

New integration test: internal/tests/core_flow_test.go

Covers: register → create household → create category → create account → upsert budget → create transaction → list budgets

Verifies that budget actual_cents reflects transactions.

internal/tests/auth_test.go: minor fixes (slogDiscard helper, io import).

Files Added

internal/security/kek.go

internal/tests/core_flow_test.go

routes/util.go

internal/db/migrations/0003_future_security_columns.sql

MIGRATION_NOTES.md (this document)

Files Modified

routes/auth.go: KEK derivation, import path fix, parameter list fix.

routes/router.go: rate limiting on auth; DELETE budget route.

Dockerfile: build target updated; removed .env copy.

cmd/server/main.go: removed unused import.

internal/tests/auth_test.go: test helpers/imports fixed.

go.mod/go.sum: cleaned (removed gorilla/mux).

Acceptance/Verification

No references to "bookkeeper-backend/models" remain.

No references to gorilla/mux remain.

go build ./... succeeds.

Auth endpoints return rate limit headers.

DELETE /v1/households/{id}/budgets/{budgetID} implemented.

Core flow integration test passes.

Docker image builds successfully.

Suggested Commands

grep -R ""bookkeeper-backend/models"" . || echo "OK: no legacy model imports"

grep -R "gorilla/mux" . || echo "OK: no gorilla/mux"

go build ./...

go test ./...

docker build -t bookkeeper-backend .

Next Steps (Future Stages)

Password verifier separation: use dedicated password_verifier column per the placeholder migration.

Refresh token reuse detection with token family rotation.

Enhanced encryption: data-at-rest using derived DEKs.

Audit logging for sensitive operations.

MFA: TOTP/WebAuthn support.

Performance: DB pooling and caching.

Monitoring: metrics and health checks.

Notes

The repository was already largely organized under /v1 with few legacy remnants.

The main work focused on KEK derivation, rate limiting, integration testing, and dependency cleanup.

Rate limiting currently uses in-memory storage suitable for single-instance deployments.

How to fix the conflict in your file

Keep the merged content above.

Remove all conflict markers: <<<<<<<, =======, >>>>>>>.

Commit the resolved file:

git add MIGRATION_NOTES.md

git commit -m "Resolve merge conflict in migration notes; consolidate stage 1 documentation"