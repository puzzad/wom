CREATE FUNCTION "public"."getteamcodejwt"("teamcode" "text") RETURNS "text"
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
DECLARE
    date      timestamp with time zone;
    after     int;
    expires   int;
    jwtsecret text;
    projectid text;
    jwt       json;
    codeMatch bool;
    response  text;
BEGIN
    SELECT NOW() into date;
    select extract(epoch from (date - interval '1 hour')) into after;
    SELECT extract(epoch from (date + interval '1 day')) INTO expires;
    SELECT value FROM internal.secrets WHERE name = 'jwtsecret' INTO jwtsecret;
    SELECT value FROM internal.secrets WHERE name = 'projectid' INTO projectid;
    SELECT EXISTS(SELECT code FROM games WHERE code = teamcode) INTO codeMatch;
    SELECT FORMAT(
                   '{"iss" : "supabase", "ref" : "%s", "role": "authenticated", "iat": %s, "exp": %s, "code": "%s"}',
                   projectid, after, expires, teamcode)
    INTO jwt;
    IF codeMatch is TRUE THEN
        SELECT extensions.sign(jwt, jwtsecret) INTO response;
        return response;
    END IF;
    RAISE EXCEPTION 'Invalid team code';
END
$$;
ALTER FUNCTION "public"."getteamcodejwt"("teamcode" "text") OWNER TO "postgres";
