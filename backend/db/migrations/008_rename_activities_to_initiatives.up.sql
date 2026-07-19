-- Migration: Rename activities to initiatives and create tasks table
-- This restructures the data model: Initiative > Task > Evidence

-- 1. Rename activities table to initiatives
ALTER TABLE activities RENAME TO initiatives;

-- 2. Rename activity_pillars to initiative_pillars
ALTER TABLE activity_pillars RENAME TO initiative_pillars;
ALTER TABLE initiative_pillars RENAME COLUMN activity_id TO initiative_id;

-- 3. Rename indexes
ALTER INDEX idx_activities_pdi RENAME TO idx_initiatives_pdi;
ALTER INDEX idx_activities_user_created RENAME TO idx_initiatives_user_created;
ALTER INDEX idx_activities_user_completed RENAME TO idx_initiatives_user_completed;
ALTER INDEX idx_activities_user_progress RENAME TO idx_initiatives_user_progress;

-- 4. Create tasks table
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    initiative_id UUID NOT NULL REFERENCES initiatives(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    execution TEXT, -- O que foi feito (campo de execucao)
    progress_percentage INTEGER NOT NULL DEFAULT 0 CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    completed_at DATE,
    created_at DATE DEFAULT CURRENT_DATE,
    updated_at DATE DEFAULT CURRENT_DATE
);

CREATE INDEX idx_tasks_initiative ON tasks(initiative_id, created_at DESC);

-- 5. Create task_evidences table (replaces activity_evidences)
CREATE TABLE task_evidences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    evidence_url TEXT NOT NULL,
    description TEXT,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE INDEX idx_task_evidences_task ON task_evidences(task_id, created_at DESC);

-- 6. Migrate existing evidences: create a default task per initiative that has evidences
-- Then move evidences to the new task_evidences table
INSERT INTO tasks (id, initiative_id, title, execution, progress_percentage, created_at)
SELECT DISTINCT
    uuid_generate_v4(),
    ae.activity_id,
    'Atividade migrada',
    'Migrada automaticamente da estrutura anterior',
    100,
    ae.created_at
FROM activity_evidences ae;

-- Move evidences to task_evidences linked to the migrated tasks
INSERT INTO task_evidences (id, task_id, evidence_url, description, created_at)
SELECT
    ae.id,
    t.id,
    ae.evidence_url,
    ae.description,
    ae.created_at
FROM activity_evidences ae
JOIN tasks t ON t.initiative_id = ae.activity_id AND t.title = 'Atividade migrada';

-- 7. Drop old activity_evidences table
DROP TABLE activity_evidences;

-- 8. Rename constraint on initiative_pillars
ALTER TABLE initiative_pillars DROP CONSTRAINT activity_pillars_pkey;
ALTER TABLE initiative_pillars ADD PRIMARY KEY (initiative_id, pillar);
