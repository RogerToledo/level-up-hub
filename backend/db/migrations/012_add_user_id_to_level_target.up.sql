-- Migration: Make level_target user-specific
-- Each user defines their own target level per year

-- 1. Drop the old unique constraint (one target per year globally)
ALTER TABLE level_target DROP CONSTRAINT level_target_unique_year;

-- 2. Add user_id column
ALTER TABLE level_target ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;

-- 3. Add unique constraint per user per year
ALTER TABLE level_target ADD CONSTRAINT level_target_unique_user_year UNIQUE (user_id, year);

-- 4. Create index for user lookups
CREATE INDEX idx_level_target_user ON level_target(user_id, year DESC);
