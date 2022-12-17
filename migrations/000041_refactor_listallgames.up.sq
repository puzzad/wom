DROP FUNCTION "public"."listallgames";
CREATE FUNCTION "public"."listallgames"()
    RETURNS TABLE
            (
                "id"            BIGINT,
                "status"        VARCHAR,
                "code"          VARCHAR,
                "user"          uuid,
                "adventure"     BIGINT,
                "puzzle"        BIGINT,
                "startTime"     TIMESTAMP WITH TIME ZONE,
                "endTime"       TIMESTAMP WITH TIME ZONE,
                "adventurename" VARCHAR,
                "puzzletitle"   TEXT,
                "useremail"     VARCHAR
            )
    LANGUAGE "plpgsql"
    SECURITY DEFINER
AS
$$
BEGIN
    IF auth.checkadmin(auth.jwt()) THEN
        RETURN QUERY SELECT "games"."id",
                            "games"."status",
                            "games"."code",
                            "games"."user",
                            "games"."adventure",
                            "games"."puzzle",
                            "games"."startTime",
                            "games"."endTime",
                            "adventures"."name",
                            "puzzles"."title",
                            "users"."email"
                     FROM "games"
                              LEFT JOIN "adventures" ON "games"."adventure" = "adventures"."id"
                              LEFT JOIN "puzzles" ON "games"."puzzle" = "puzzles"."id"
                              LEFT JOIN "auth"."users" ON "games"."user" = "auth"."users"."id";
    END IF;
END
$$;
ALTER FUNCTION "public"."listallgames"() OWNER TO "postgres";
