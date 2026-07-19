-- Rollback: Revert initiatives back to activities

-- 1. Recreate activity_evidences from task_evidences
CREATE TABLE activity_evidences (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    activity_id UUID NOT NULL,
    evidence_url TEXT NOT NULL,
    description TEXT,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    CONSTRAINT activity_evidences_pk PRIMARY KEY (id)
);

-- Move evidences back (link through tasks to initiatives)
INSERT INTO activity_evidences (id, activity_id, evidence_url, description, created_at)
SELECT
    te.id,
    t.initiative_id,
    te.evidence_url,
    te.description,
    te.created_at
FROM task_evidences te
JOIN tasks t ON t.id = te.task_id;

-- 2. Drop new tables
DROP TABLE task_evidences;
DROP TABLE tasks;

-- 3. Rename initiative_pillars back
ALTER TABLE initiative_pillars DROP CONSTRAINT initiative_pillars_pkey;
ALTER TABLE initiative_pillars RENAME COLUMN initiative_id TO activity_id;
ALTER TABLE initiative_pillars RENAME TO activity_pillars;
ALTER TABLE activity_pillars ADD PRIMARY KEY (activity_id, pillar);

-- 4. Rename indexes back
ALTER INDEX idx_initiatives_pdi RENAME TO idx_activities_pdi;
ALTER INDEX idx_initiatives_user_created RENAME TO idx_activities_user_created;
ALTER INDEX idx_initiatives_user_completed RENAME TO idx_activities_user_completed;
ALTER INDEX idx_initiatives_user_progress RENAME TO idx_activities_user_progress;

-- 5. Rename initiatives back to activities
ALTER TABLE initiatives RENAME TO activities;

-- 6. Add back FK on activity_evidences
ALTER TABLE activity_evidences
ADD CONSTRAINT activity_evidences_activity_fk 
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE;
