INSERT INTO "storage"."buckets" (id, name, public)
VALUES ('extras', 'extras', true)
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS public.extras
(
    id          serial primary key,
    filename    varchar unique not null,
    title       varchar,
    description varchar,
    sort        int
);