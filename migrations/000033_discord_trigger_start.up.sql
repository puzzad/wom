CREATE OR REPLACE FUNCTION "public"."startadventure"() RETURNS bigint
    LANGUAGE "plpgsql" SECURITY DEFINER
AS $$DECLARE
    res int;
BEGIN
    UPDATE games SET status = 'ACTIVE', puzzle = (SELECT adventures."firstPuzzle" FROM adventures WHERE adventures.id = games.adventure), "startTime" = NOW() WHERE code = auth.gameCode() AND status = 'PAID' AND puzzle IS NULL;
    EXECUTE internal.announce(format(':checkered_flag: `%s` started', auth.gameCode()));
    SELECT puzzle FROM games WHERE code = auth.gameCode() INTO res;
    return res;
END$$;