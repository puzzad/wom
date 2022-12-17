DROP FUNCTION IF EXISTS "public"."listallusers";
CREATE FUNCTION "public"."listallusers"()
    RETURNS TABLE
            (
                "id"         uuid,
                "email"      VARCHAR,
                "created"    TIMESTAMP WITH TIME ZONE,
                "confirmed"  TIMESTAMP WITH TIME ZONE,
                "lastSignin" TIMESTAMP WITH TIME ZONE,
                "adventures" JSON
            )
    LANGUAGE "plpgsql"
    SECURITY DEFINER
AS
$$
BEGIN
    IF auth.checkadmin(auth.jwt()) THEN
        RETURN QUERY SELECT "users"."id",
                            "users"."email",
                            "users"."created_at",
                            "users"."confirmed_at",
                            "users"."last_sign_in_at",
                            COALESCE((SELECT JSON_AGG(JSON_BUILD_OBJECT('code', "code", 'adventure', "adventures"."name"))
                                      FROM "games"
                                      LEFT JOIN "adventures" on "games"."adventure" = "adventures"."id"
                                      WHERE "games"."user" = "users"."id"), '[]'::json) AS adventures
                     FROM "auth"."users";
    END IF;
END
$$;
ALTER FUNCTION "public"."listallusers"() OWNER TO "postgres";
