DO $$ BEGIN
    CREATE TYPE ladder_level AS ENUM ('P1', 'P2', 'P3', 'LT1', 'LT2', 'LT3', 'LT4');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE pillar AS ENUM ('TECHNICAL', 'RESULTS', 'INFLUENCE');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('user', 'admin');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    role user_role NOT NULL DEFAULT 'user',
    current_level public.ladder_level NOT NULL DEFAULT 'P1'::ladder_level,
    created_at DATE NOT NULL DEFAULT NOW(),
    updated_at DATE NOT NULL DEFAULT NOW()
);
