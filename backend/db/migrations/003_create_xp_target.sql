CREATE TABLE public.xp_target (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    target int4 NOT NULL,
    year int NOT NULL,
    ladder_id uuid NOT NULL,
    CONSTRAINT xp_target_pk PRIMARY KEY (id),
    CONSTRAINT xp_target_career_ladder_fk FOREIGN KEY (ladder_id) REFERENCES public.career_ladder(id),
    CONSTRAINT xp_target_unique_level_year UNIQUE (ladder_id, year)
);