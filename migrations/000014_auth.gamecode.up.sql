CREATE FUNCTION "auth"."gamecode"() RETURNS character varying
    LANGUAGE "plpgsql" STABLE
    AS $$BEGIN
return auth.jwt()->>'code';
END;$$;
ALTER FUNCTION "auth"."gamecode"() OWNER TO "postgres";
