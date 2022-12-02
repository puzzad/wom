CREATE FUNCTION "internal"."announce"("message" character varying) RETURNS "void"
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
declare
    url     VARCHAR;
    payload jsonb;
BEGIN
    SELECT value FROM internal.secrets WHERE name = 'discord-events-webhook' INTO url;
    payload := json_build_object(
            'content', message
        );
    perform net.http_post(url := url, body := payload);
END;
$$;
ALTER FUNCTION "internal"."announce"("message" character varying) OWNER TO "postgres";
