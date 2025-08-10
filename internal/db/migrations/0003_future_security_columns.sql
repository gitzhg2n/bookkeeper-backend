-- Migration: 0003_future_security_columns.sql
-- Purpose: Placeholder for future password security enhancements
-- Status: Non-breaking, idempotent
-- Stage: Future (not Stage 1)

-- This migration is a placeholder for future security column additions
-- such as password_verifier separation from KEK derivation.
-- Currently, no changes are made to maintain backward compatibility.

-- Future columns to consider (commented out for now):
-- ALTER TABLE users ADD COLUMN password_verifier BLOB;
-- ALTER TABLE users ADD COLUMN kek_version INTEGER DEFAULT 1;
-- ALTER TABLE users ADD COLUMN security_flags INTEGER DEFAULT 0;

-- For now, this migration serves as a version marker
-- and ensures the migration system recognizes this step.

-- No-op statement to make this a valid migration
SELECT 1 as placeholder_migration_0003;