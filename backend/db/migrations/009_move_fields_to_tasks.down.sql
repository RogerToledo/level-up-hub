-- Rollback: Move fields back from tasks to initiatives

-- 1. Re-add ladder_id and impact_summary to initiatives
ALTER TABLE initiatives
ADD COLUMN ladder_id UUID REFERENCES career_ladder(id),
ADD COLUMN impact_summary TEXT;

-- 2. Copy ladder_id back from first task of each initiative
UPDATE initiatives i
SET ladder_id = (
    SELECT t.ladder_id FROM tasks t WHERE t.initiative_id = i.id LIMIT 1
);

-- 3. Recreate initiative_pillars from task_pillars
CREATE TABLE initiative_pillars (
    initiative_id UUID NOT NULL REFERENCES initiatives(id) ON DELETE CASCADE,
    pillar pillar NOT NULL,
    PRIMARY KEY (initiative_id, pillar)
);

INSERT INTO initiative_pillars (initiative_id, pillar)
SELECT DISTINCT t.initiative_id, tp.pillar
FROM task_pillars tp
JOIN tasks t ON t.id = tp.task_id;

-- 4. Drop task_pillars
DROP TABLE task_pillars;

-- 5. Remove ladder_id and impact_summary from tasks
ALTER TABLE tasks
DROP COLUMN ladder_id,
DROP COLUMN impact_summary;
