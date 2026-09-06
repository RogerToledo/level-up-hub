ALTER TABLE plans ADD COLUMN position INTEGER NOT NULL DEFAULT 0;

UPDATE plans SET position = (
    SELECT COUNT(*) FROM plans p2
    WHERE p2.user_id = plans.user_id AND p2.created_at <= plans.created_at
) - 1;

CREATE INDEX idx_plans_user_position ON plans(user_id, position);