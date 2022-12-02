CREATE FUNCTION "public"."getstats"("gamecode" character varying) RETURNS TABLE("title" "text", "solvetime" timestamp with time zone, "hints" bigint)
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
BEGIN
    if auth.gamecode() <> gameCode then
        raise exception 'Invalid game code provided';
    end if;
    RETURN QUERY
        SELECT p.title                                AS title,
               (SELECT MAX(g.created_at)
                FROM guesses g
                WHERE g.game = gameCode
                  AND g.puzzle = p.id
                  AND g.correct)                      AS solveTime,
               (SELECT COUNT(*)
                FROM internal.usedhints u
                         JOIN internal.hints h on h.id = u.hint
                WHERE h.puzzle = p.id
                  AND u.game = gameCode) AS hints
        FROM puzzles p
        WHERE adventure = (SELECT adventure FROM games WHERE code = gameCode)
        ORDER BY solveTime;
END;
$$;
ALTER FUNCTION "public"."getstats"("gamecode" character varying) OWNER TO "postgres";
