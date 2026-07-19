-- Migration: Rename xp_target to level_target

ALTER TABLE xp_target RENAME TO level_target;

ALTER TABLE level_target RENAME CONSTRAINT xp_target_pk TO level_target_pk;
ALTER TABLE level_target RENAME CONSTRAINT xp_target_career_ladder_fk TO level_target_career_ladder_fk;
ALTER TABLE level_target RENAME CONSTRAINT xp_target_unique_level_year TO level_target_unique_level_year;
