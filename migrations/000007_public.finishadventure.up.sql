CREATE FUNCTION "public"."finishadventure"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$BEGIN
  if NEW.puzzle is null then
    NEW.status = 'EXPIRED';
    NEW."endTime" = NOW();
  end if;
  return NEW;
END;
$$;
ALTER FUNCTION "public"."finishadventure"() OWNER TO "postgres";
CREATE TRIGGER "FinishAdventure" BEFORE UPDATE ON "public"."games" FOR EACH ROW EXECUTE FUNCTION "public"."finishadventure"();
