CREATE FUNCTION "public"."checkguess"() RETURNS "trigger"
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$BEGIN
  NEW.correct = (select count(*) > 0 FROM internal.answers WHERE puzzle = NEW.puzzle AND lower(answer) = lower(NEW.content));
  RETURN NEW;
END;$$;
ALTER FUNCTION "public"."checkguess"() OWNER TO "postgres";
CREATE TRIGGER "CheckGuess" BEFORE INSERT ON "public"."guesses" FOR EACH ROW EXECUTE FUNCTION "public"."checkguess"();
