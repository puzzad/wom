CREATE TABLE IF NOT EXISTS public.content
(
    page    VARCHAR NOT NULL PRIMARY KEY,
    content TEXT
);

ALTER TABLE "public"."content" OWNER TO "postgres";
ALTER TABLE "public"."content" ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "Read access to content for everyone" ON "public"."content";
CREATE POLICY "Read access to content for everyone" ON "public"."content" FOR SELECT USING (true);

INSERT INTO public.content(page, content)
VALUES ('homepage',
        '<p>Welcome to puzzad! This is a placeholder.</p>' ||
        '<p>The site admin should update this text in the `public.content` table.</p>'),
       ('privacy',
        '<p>This is a placeholder.</p>' ||
        '<p>The site admin should update this text in the `public.content` table.</p>')
ON CONFLICT DO NOTHING;