-- Rollback initial schema migration

DROP INDEX IF EXISTS idx_households_created_by;
DROP TABLE IF EXISTS households;

DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP TABLE IF EXISTS refresh_tokens;

DROP TABLE IF EXISTS users;
