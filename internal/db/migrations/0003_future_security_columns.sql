 copilot/fix-184f7982-e511-4e6f-9dc2-305d1c6b4c15
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

-- 0003_future_security_columns.sql
-- Migration 0003: Future Security Columns (Placeholder)
-- Purpose: Reserve migration slot for future password verification separation.
-- This migration is idempotent and non-breaking; it performs no schema changes today.

-- Planned future changes (uncomment/adjust when implementing):
-- ALTER TABLE users ADD COLUMN password_verifier BLOB; -- for password authentication only
-- ALTER TABLE users ADD COLUMN verifier_salt BLOB; -- salt used for verifier derivation
-- ALTER TABLE users ADD COLUMN verifier_version INTEGER DEFAULT 1;
-- Optional transition fields:
-- ALTER TABLE users ADD COLUMN kek_version INTEGER DEFAULT 1; -- track KEK derivation version

-- Rationale:
-- 1) password_verifier: decouple authentication from KEK/DEK derivation
-- 2) Maintain existing password_hash temporarily for KEK derivation while migrating
-- 3) Enable versioned upgrades of verifier/KEK derivation

-- No-op to keep migration sequence consistent
SELECT 1 AS placeholder;

How to fix the conflict in your file

Replace the entire file content with the merged version above.

Remove all conflict markers: <<<<<<<, =======, >>>>>>>.

Commit the resolution:

git add internal/db/migrations/0003_future_security_columns.sql

git commit -m "Resolve merge conflict in migration 0003; add unified placeholder with future plan comments"
 main
