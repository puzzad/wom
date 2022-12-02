CREATE OR REPLACE FUNCTION  "auth"."checkadmin"(jwt jsonb) RETURNS bool
    LANGUAGE "plpgsql"
    SECURITY DEFINER
AS
$$
DECLARE
    admindomain  text;
    domainlength int;
BEGIN
    SELECT value FROM "internal"."secrets" where name = 'admindomain' INTO admindomain;
    SELECT LENGTH(admindomain) INTO domainlength;

    RETURN right(jwt ->> 'email', domainlength) = admindomain;
END
$$;
ALTER FUNCTION "auth"."checkadmin"("jwt" "jsonb") OWNER TO "postgres";

DROP POLICY IF EXISTS "Admins can view auth.users" ON "auth"."users";
CREATE POLICY "Admins can view auth.users"
    ON auth.users
    FOR SELECT USING (
    auth.checkadmin(auth.jwt())
    );
