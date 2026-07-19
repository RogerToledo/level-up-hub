-- Rollback: Revert level_target back to int target + ladder_id

ALTER TABLE level_target DROP CONSTRAINT level_target_unique_year;
ALTER TABLE level_target DROP COLUMN target;

ALTER TABLE level_target ADD COLUMN target int4 NOT NULL DEFAULT 0;
ALTER TABLE level_target ADD COLUMN ladder_id uuid;

ALTER TABLE level_target ADD CONSTRAINT level_target_career_ladder_fk
    FOREIGN KEY (ladder_id) REFERENCES public.career_ladder(id);
ALTER TABLE level_target ADD CONSTRAINT level_target_unique_level_year
    UNIQUE (ladder_id, year);
