CREATE FUNCTION "public"."startadventure"() RETURNS bigint
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$DECLARE
  res int;
BEGIN
  UPDATE games SET status = 'ACTIVE', puzzle = (SELECT adventures."firstPuzzle" FROM adventures WHERE adventures.id = games.adventure), "startTime" = NOW() WHERE code = auth.gameCode() AND status = 'PAID' AND puzzle IS NULL;
  SELECT puzzle FROM games WHERE code = auth.gameCode() INTO res;
  return res;
END$$;
ALTER FUNCTION "public"."startadventure"() OWNER TO "postgres";
