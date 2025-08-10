-- Migration 0003: Future Security Columns (Placeholder)
-- This migration is a placeholder for future security improvements.
-- It will add columns for enhanced password verification in a future stage.

-- Uncomment and modify as needed in future iterations:
-- ALTER TABLE users ADD COLUMN password_verifier TEXT;
-- ALTER TABLE users ADD COLUMN kek_version INTEGER DEFAULT 1;

-- For now, this migration does nothing (idempotent)
SELECT 1;