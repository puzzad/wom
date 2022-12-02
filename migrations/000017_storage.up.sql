INSERT INTO "storage"."buckets" (id, name, public) VALUES ('puzzles', 'puzzles', false);
INSERT INTO "storage"."buckets" (id, name, public) VALUES ('adventures', 'adventures', true);

ALTER TABLE "storage"."buckets" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "storage"."migrations" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "storage"."objects" ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Read access for users with code" ON "storage"."objects" FOR SELECT USING ((("bucket_id" = 'puzzles'::"text") AND (("storage"."foldername"("name"))[1] IN ( SELECT "p"."storage_slug"
   FROM ("public"."puzzles" "p"
     JOIN "public"."games" "g" ON ((("g"."puzzle" = "p"."id") AND (("g"."code")::"text" = ("auth"."gamecode"())::"text"))))))));
