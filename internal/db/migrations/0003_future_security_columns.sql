-- 0003_future_security_columns.sql
-- Placeholder migration for future password verifier column separation
-- This migration is idempotent and non-breaking for current functionality

-- Future: Add separate password_verifier column when ready to separate 
-- password verification from DEK derivation key material
-- 
-- Planned changes:
-- ALTER TABLE users ADD COLUMN password_verifier BLOB;
-- ALTER TABLE users ADD COLUMN verifier_salt BLOB;
-- ALTER TABLE users ADD COLUMN verifier_version INTEGER DEFAULT 1;
--
-- This will enable:
-- 1. password_verifier: For password authentication only
-- 2. password_hash: Retained for KEK derivation during transition
-- 3. Eventual migration to separate key material completely

-- No-op migration for now - ensures migration numbering consistency
SELECT 1 AS placeholder;