-- Generated with pgAdmin

-- Table: public.snippets

-- DROP TABLE public.snippets;

CREATE TABLE IF NOT EXISTS public.snippets
(
    id SERIAL NOT NULL PRIMARY KEY,
    solution_id bigint,
    text text COLLATE pg_catalog."default",
    language text COLLATE pg_catalog."default"
)

TABLESPACE pg_default;

ALTER TABLE public.snippets
    OWNER to postgres;

COMMENT ON COLUMN public.snippets.solution_id
    IS 'TO BE FIXED: must be a foreign key to solution table!';