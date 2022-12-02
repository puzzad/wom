ALTER TABLE "internal"."adjectives" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "internal"."animals" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "internal"."answers" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "internal"."colours" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "internal"."hints" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "internal"."secrets" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "internal"."usedhints" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "public"."adventures" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "public"."games" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "public"."guesses" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "public"."puzzles" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "storage"."buckets" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "storage"."migrations" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "storage"."objects" ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Read access for users with code" ON "public"."puzzles" FOR SELECT USING (("id" = ( SELECT "games"."puzzle"
   FROM "public"."games"
  WHERE (("games"."code")::"text" = ("auth"."gamecode"())::"text"))));
CREATE POLICY "Read access for users with the code" ON "public"."guesses" FOR SELECT USING ((("game")::"text" = ("auth"."gamecode"())::"text"));
CREATE POLICY "Read access to adventures for authenticated users" ON "public"."adventures" FOR SELECT USING (("public" = true));
CREATE POLICY "Read access to games via code" ON "public"."games" FOR SELECT USING ((("code")::"text" = ("auth"."gamecode"())::"text"));
CREATE POLICY "Read access to own games" ON "public"."games" FOR SELECT USING (("auth"."uid"() = "user"));
CREATE POLICY "Read access to users with game codes" ON "public"."adventures" FOR SELECT USING (("auth"."gamecode"() IS NOT NULL));
CREATE POLICY "Write access for users with codes" ON "public"."guesses" FOR INSERT WITH CHECK ((("game")::"text" = ("auth"."gamecode"())::"text"));
CREATE POLICY "Read access for users" ON "storage"."objects" FOR SELECT USING ((("bucket_id" = 'puzzles'::"text") AND (("storage"."foldername"("name"))[1] IN ( SELECT "p"."storage_slug"
   FROM ("public"."puzzles" "p"
     JOIN "public"."games" "g" ON ((("g"."puzzle" = "p"."id") AND (("g"."code")::"text" = ("auth"."gamecode"())::"text"))))))));
