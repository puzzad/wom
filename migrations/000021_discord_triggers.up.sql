CREATE OR REPLACE FUNCTION "internal"."announce_new_user"() RETURNS "trigger"
    LANGUAGE "plpgsql" SECURITY DEFINER
AS $$BEGIN
    EXECUTE internal.announce(format('New user signup: `%s`', NEW.id));
    return null;
END;$$;

ALTER FUNCTION "internal"."announce_new_user"() OWNER TO "postgres";

CREATE OR REPLACE TRIGGER "announce_new_user" AFTER INSERT ON "auth"."users" FOR EACH ROW EXECUTE FUNCTION "internal"."announce_new_user"();



CREATE OR REPLACE FUNCTION "internal"."announce_new_game"() RETURNS "trigger"
    LANGUAGE "plpgsql" SECURITY DEFINER
AS $$DECLARE
    adventure varchar;
BEGIN
    SELECT name FROM public.adventures WHERE id = NEW.adventure INTO adventure;
    EXECUTE internal.announce(format('New game created: `%s` (adventure: %s)', NEW.code, adventure));
    return null;
END;$$;

ALTER FUNCTION "internal"."announce_new_game"() OWNER TO "postgres";

CREATE OR REPLACE TRIGGER "announce_new_game" AFTER INSERT ON "public"."games" FOR EACH ROW EXECUTE FUNCTION "internal"."announce_new_game"();



CREATE OR REPLACE FUNCTION "internal"."announce_new_guess"() RETURNS "trigger"
    LANGUAGE "plpgsql" SECURITY DEFINER
AS $$DECLARE
    puzzle varchar;
BEGIN
    SELECT title FROM public.puzzles WHERE id = NEW.puzzle INTO puzzle;
    if NEW.correct then
        EXECUTE internal.announce(format(':tada: `%s`/`%s`: `%s`', NEW.game, puzzle, NEW.content));
    else
        EXECUTE internal.announce(format(':x: `%s`/`%s`: `%s`', NEW.game, puzzle, NEW.content));
    end if;
    return null;
END;$$;

ALTER FUNCTION "internal"."announce_new_guess"() OWNER TO "postgres";

CREATE OR REPLACE TRIGGER "announce_new_guess" AFTER INSERT ON "public"."guesses" FOR EACH ROW EXECUTE FUNCTION "internal"."announce_new_guess"();
