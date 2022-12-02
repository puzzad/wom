CREATE OR REPLACE FUNCTION "internal"."announce_new_guess"() RETURNS "trigger"
    LANGUAGE "plpgsql" SECURITY DEFINER
AS $$DECLARE
    puzzle varchar;
BEGIN
    SELECT title FROM public.puzzles WHERE id = NEW.puzzle INTO puzzle;
    if NEW.content = '*hint' then
        EXECUTE internal.announce(format(':bulb: `%s`/`%s`: hint requested', NEW.game, puzzle));
    elseif NEW.correct then
        EXECUTE internal.announce(format(':tada: `%s`/`%s`: `%s`', NEW.game, puzzle, NEW.content));
    else
        EXECUTE internal.announce(format(':x: `%s`/`%s`: `%s`', NEW.game, puzzle, NEW.content));
    end if;
    return null;
END;$$;
