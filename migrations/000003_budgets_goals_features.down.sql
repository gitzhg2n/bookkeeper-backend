-- Rollback budgets, goals and features

DROP INDEX IF EXISTS idx_bills_is_active;
DROP INDEX IF EXISTS idx_bills_next_due_date;
DROP INDEX IF EXISTS idx_bills_household_id;
DROP TABLE IF EXISTS bills;

DROP INDEX IF EXISTS idx_notifications_type;
DROP INDEX IF EXISTS idx_notifications_read_at;
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP TABLE IF EXISTS notifications;

DROP TABLE IF EXISTS user_settings;

DROP INDEX IF EXISTS idx_goals_due_date;
DROP INDEX IF EXISTS idx_goals_user_id;
DROP TABLE IF EXISTS goals;

DROP INDEX IF EXISTS idx_budgets_category_id;
DROP INDEX IF EXISTS idx_budgets_month;
DROP INDEX IF EXISTS idx_budgets_household_id;
DROP TABLE IF EXISTS budgets;
