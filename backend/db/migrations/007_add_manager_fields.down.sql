DROP INDEX IF EXISTS idx_users_manager_email;

ALTER TABLE users
DROP COLUMN IF EXISTS manager_name,
DROP COLUMN IF EXISTS manager_email;
