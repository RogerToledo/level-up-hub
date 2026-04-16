CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ladder_id UUID NOT NULL REFERENCES career_ladder(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    is_pdi_target BOOLEAN NOT NULL DEFAULT FALSE,
    progress_percentage INTEGER NOT NULL DEFAULT 0 CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    impact_summary TEXT, 
    completed_at DATE,
    created_at DATE DEFAULT CURRENT_DATE,
    updated_at DATE DEFAULT CURRENT_DATE
);

CREATE INDEX idx_activities_pdi ON activities(user_id, is_pdi_target);

CREATE TABLE activity_pillars (
    activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    pillar pillar NOT NULL,
    PRIMARY KEY (activity_id, pillar)
);

CREATE TABLE activity_evidences (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    activity_id uuid NOT NULL,
    evidence_url TEXT NOT NULL,
    description TEXT, -- Ex: "Migration PR", "Slack Screenshot"
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    CONSTRAINT activity_evidences_pk PRIMARY KEY (id),
    CONSTRAINT activity_evidences_activity_fk FOREIGN KEY (activity_id) 
        REFERENCES public.activities(id) ON DELETE CASCADE
);