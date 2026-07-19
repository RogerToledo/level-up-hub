-- Rollback: Remove user_id from level_target

DROP INDEX IF EXISTS idx_level_target_user;
ALTER TABLE level_target DROP CONSTRAINT level_target_unique_user_year;
ALTER TABLE level_target DROP COLUMN user_id;
ALTER TABLE level_target ADD CONSTRAINT level_target_unique_year UNIQUE (year);
