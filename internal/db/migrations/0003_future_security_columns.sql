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