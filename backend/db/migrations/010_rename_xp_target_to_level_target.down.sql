-- Rollback: Rename level_target back to xp_target

ALTER TABLE level_target RENAME CONSTRAINT level_target_pk TO xp_target_pk;
ALTER TABLE level_target RENAME CONSTRAINT level_target_career_ladder_fk TO xp_target_career_ladder_fk;
ALTER TABLE level_target RENAME CONSTRAINT level_target_unique_level_year TO xp_target_unique_level_year;

ALTER TABLE level_target RENAME TO xp_target;
