-- Migration: Add is_extra field to tasks
-- Marks tasks that are overdelivery (extra) within an initiative

ALTER TABLE tasks ADD COLUMN is_extra BOOLEAN NOT NULL DEFAULT FALSE;
