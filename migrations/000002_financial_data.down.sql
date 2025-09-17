-- Rollback financial data tables

DROP INDEX IF EXISTS idx_transactions_occurred_at;
DROP INDEX IF EXISTS idx_transactions_category_id;
DROP INDEX IF EXISTS idx_transactions_user_id;
DROP INDEX IF EXISTS idx_transactions_account_id;
DROP TABLE IF EXISTS transactions;

DROP INDEX IF EXISTS idx_categories_parent_id;
DROP INDEX IF EXISTS idx_categories_household_id;
DROP TABLE IF EXISTS categories;

DROP INDEX IF EXISTS idx_accounts_archived_at;
DROP INDEX IF EXISTS idx_accounts_type;
DROP INDEX IF EXISTS idx_accounts_household_id;
DROP TABLE IF EXISTS accounts;

DROP INDEX IF EXISTS idx_household_members_user_id;
DROP INDEX IF EXISTS idx_household_members_household_id;
DROP TABLE IF EXISTS household_members;
