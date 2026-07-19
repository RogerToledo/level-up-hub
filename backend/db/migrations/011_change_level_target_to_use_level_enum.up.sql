-- Migration: Change level_target.target from int to ladder_level enum
-- Remove ladder_id (redundant), target becomes the level name (P1, P2, etc.)

-- 1. Drop old constraints
ALTER TABLE level_target DROP CONSTRAINT level_target_career_ladder_fk;
ALTER TABLE level_target DROP CONSTRAINT level_target_unique_level_year;

-- 2. Remove old columns
ALTER TABLE level_target DROP COLUMN target;
ALTER TABLE level_target DROP COLUMN ladder_id;

-- 3. Add new target column as ladder_level enum
ALTER TABLE level_target ADD COLUMN target ladder_level NOT NULL DEFAULT 'P1';

-- 4. Add unique constraint (one target per year)
ALTER TABLE level_target ADD CONSTRAINT level_target_unique_year UNIQUE (year);
