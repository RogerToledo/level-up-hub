CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    initiative_id UUID REFERENCES initiatives(id) ON DELETE SET NULL,
    level_target TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    created_at DATE NOT NULL DEFAULT NOW(),
    updated_at DATE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_plans_user_id ON plans(user_id);
CREATE INDEX idx_plans_initiative_id ON plans(initiative_id);
CREATE INDEX idx_plans_status ON plans(status);