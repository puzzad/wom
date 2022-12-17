CREATE OR REPLACE FUNCTION "public"."listallgames"(jwt TEXT)
    RETURNS TABLE
            (
                "id"        BIGINT,
                "status"    VARCHAR,
                "code"      VARCHAR,
                "user"      uuid,
                "adventure" BIGINT,
                "puzzle"    BIGINT,
                "startTime" TIMESTAMP WITH TIME ZONE,
                "endTime"   TIMESTAMP WITH TIME ZONE
            )
    LANGUAGE "plpgsql"
    SECURITY DEFINER
AS
$$
DECLARE
    verified BOOL;
    result TEXT;
BEGIN
    SELECT payload ->> 'role', valid::BOOL FROM extensions.verify(jwt, (SELECT value FROM internal.secrets WHERE name = 'jwtsecret')) INTO result, verified;
    IF verified AND result = 'supabase_admin' THEN
        RETURN QUERY SELECT "games"."id",
                            "games"."status",
                            "games"."code",
                            "games"."user",
                            "games"."adventure",
                            "games"."puzzle",
                            "games"."startTime",
                            "games"."endTime"
                     FROM "games";
    END IF;
END
$$;
ALTER FUNCTION "public"."listallgames"(TEXT) OWNER TO "postgres";
