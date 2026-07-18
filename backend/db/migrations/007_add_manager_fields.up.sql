-- Migration: Add manager fields to users table
-- This allows users to register their engineering manager for report sharing

ALTER TABLE users 
ADD COLUMN manager_name TEXT,
ADD COLUMN manager_email TEXT;

-- Add index for manager email lookups
CREATE INDEX idx_users_manager_email ON users(manager_email) WHERE manager_email IS NOT NULL;

-- Add comments for documentation
COMMENT ON COLUMN users.manager_name IS 'Name of the user''s engineering manager';
COMMENT ON COLUMN users.manager_email IS 'Email address of the engineering manager for report distribution';
