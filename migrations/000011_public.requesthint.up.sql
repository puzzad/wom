CREATE FUNCTION "public"."requesthint"("puzzleid" bigint, "gamecode" character varying, "hintid" bigint) RETURNS "void"
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
BEGIN
    if auth.gamecode() <> gameCode then
        raise exception 'Invalid game code provided';
    end if;
    if not exists(select 1 from games where code = gameCode and puzzle = puzzleId) then
        raise exception 'Invalid puzzle provided';
    end if;
    if not exists(select 1 from internal.hints where id = hintId and puzzle = puzzleId) then
        raise exception 'Invalid hint provided';
    end if;
    INSERT INTO internal.usedhints (hint, game) VALUES (hintId, gameCode) ON CONFLICT DO NOTHING;
    INSERT INTO guesses (content, puzzle, game) VALUES ('*hint', puzzleId, gameCode);
END;
$$;
ALTER FUNCTION "public"."requesthint"("puzzleid" bigint, "gamecode" character varying, "hintid" bigint) OWNER TO "postgres";
