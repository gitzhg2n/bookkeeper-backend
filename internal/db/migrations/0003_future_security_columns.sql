-- Migration: 0003_future_security_columns.sql
-- Purpose: Placeholder for future password/security column separation.
-- This is currently a no-op and provides a reserved slot for:
--   - password_verifier
--   - verifier_salt
--   - verifier_version
--   - kek_version
-- Idempotent.

SELECT 1 AS placeholder_0003;