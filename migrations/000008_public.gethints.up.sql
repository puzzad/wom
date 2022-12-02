CREATE FUNCTION "public"."gethints"("puzzleid" bigint, "gamecode" character varying) RETURNS TABLE("id" bigint, "title" character varying, "locked" boolean, "message" character varying)
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
BEGIN
    if auth.gamecode() <> gameCode then
        raise exception 'Invalid game code provided';
    end if;
    if not exists(select 1 from games where code = gameCode and puzzle = puzzleId) then
        raise exception 'Invalid puzzle provided';
    end if;
    RETURN QUERY
        SELECT hints.id,
               hints.title,
               u.id is null                                          AS locked,
               case when u.id is null then '' else hints.message end AS message
        FROM internal.hints
                 LEFT JOIN internal.usedhints u on hints.id = u.hint AND u.game = gameCode
        WHERE hints.puzzle = puzzleId
        ORDER BY hints."order";
END;
$$;
ALTER FUNCTION "public"."gethints"("puzzleid" bigint, "gamecode" character varying) OWNER TO "postgres";
