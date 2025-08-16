-- Seed minimal data for E2E tests
-- Creates a test user, household, and a couple of goals and an account

BEGIN;

-- Insert a test user (avoid specifying id so DB assigns it)
INSERT INTO users (email, password_hash, created_at)
VALUES ('test@example.com', 'testhash', NOW())
ON CONFLICT (email) DO NOTHING;

-- Ensure a household exists for that user
WITH u AS (SELECT id FROM users WHERE email = 'test@example.com' LIMIT 1)
INSERT INTO households (name, created_by, created_at)
SELECT 'Test Household', u.id, NOW() FROM u
WHERE NOT EXISTS (SELECT 1 FROM households WHERE name = 'Test Household');

-- Add household member linking the test user to the household
WITH u AS (SELECT id FROM users WHERE email = 'test@example.com' LIMIT 1),
	 h AS (SELECT id FROM households WHERE name = 'Test Household' LIMIT 1)
INSERT INTO household_members (household_id, user_id, role, created_at)
SELECT h.id, u.id, 'owner', NOW() FROM h, u
WHERE NOT EXISTS (SELECT 1 FROM household_members WHERE household_id = h.id AND user_id = u.id);

-- Add a checking account for the household
WITH h AS (SELECT id FROM households WHERE name = 'Test Household' LIMIT 1)
INSERT INTO accounts (household_id, name, type, opening_balance_cents, created_at)
SELECT h.id, 'Checking', 'checking', 500000, NOW() FROM h
WHERE NOT EXISTS (SELECT 1 FROM accounts WHERE household_id = h.id AND name = 'Checking');

-- Add a simple goal associated with the test user
WITH u AS (SELECT id FROM users WHERE email = 'test@example.com' LIMIT 1)
INSERT INTO goals (user_id, name, target_cents, current_cents, due_date, created_at)
SELECT u.id, 'Car', 1000000, 250000, NOW() + INTERVAL '90 days', NOW() FROM u
WHERE NOT EXISTS (SELECT 1 FROM goals WHERE user_id = u.id AND name = 'Car');

COMMIT;
