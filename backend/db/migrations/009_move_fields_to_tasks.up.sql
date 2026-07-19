-- Migration: Move ladder_id, pillars, impact_summary from initiatives to tasks
-- Initiative now only has: title, description, is_pdi_target
-- Task now has: title, execution, ladder_id, pillars, impact_summary, progress_percentage

-- 1. Add ladder_id and impact_summary to tasks
ALTER TABLE tasks
ADD COLUMN ladder_id UUID REFERENCES career_ladder(id),
ADD COLUMN impact_summary TEXT;

-- 2. Copy ladder_id from initiative to its tasks (for existing data)
UPDATE tasks t
SET ladder_id = i.ladder_id
FROM initiatives i
WHERE t.initiative_id = i.id;

-- 3. Make ladder_id NOT NULL after backfill
ALTER TABLE tasks ALTER COLUMN ladder_id SET NOT NULL;

-- 4. Create task_pillars table
CREATE TABLE task_pillars (
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    pillar pillar NOT NULL,
    PRIMARY KEY (task_id, pillar)
);

-- 5. Copy pillars from initiative_pillars to task_pillars (for existing tasks)
INSERT INTO task_pillars (task_id, pillar)
SELECT t.id, ip.pillar
FROM tasks t
JOIN initiative_pillars ip ON ip.initiative_id = t.initiative_id;

-- 6. Drop initiative_pillars table (pillars now live on tasks)
DROP TABLE initiative_pillars;

-- 7. Remove ladder_id and impact_summary from initiatives (no longer needed)
ALTER TABLE initiatives
DROP COLUMN ladder_id,
DROP COLUMN impact_summary;
