CREATE TYPE ladder_level AS ENUM ('P1', 'P2', 'P3', 'LT1', 'LT2', 'LT3', 'LT4');
CREATE TYPE pillar AS ENUM ('TECHNICAL', 'RESULTS', 'INFLUENCE');

CREATE TABLE career_ladder (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    level ladder_level NOT NULL,
    xp_reward INTEGER NOT NULL,
    technical TEXT NOT NULL,
    expected_results TEXT NOT NULL,
    leadership_scope TEXT NOT NULL
);
