CREATE FUNCTION "public"."advancetonextpuzzle"() RETURNS "trigger"
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$
BEGIN
  if NEW.correct then
    UPDATE games SET puzzle = (SELECT next FROM puzzles WHERE id = NEW.puzzle) WHERE games.code = NEW.game;
    return NEW;
  end if;
RETURN NULL;
END;
$$;
ALTER FUNCTION "public"."advancetonextpuzzle"() OWNER TO "postgres";
CREATE TRIGGER "AdvanceToNextPuzzle" AFTER INSERT ON "public"."guesses" FOR EACH ROW EXECUTE FUNCTION "public"."advancetonextpuzzle"();
