CREATE FUNCTION "internal"."setteamcode"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$BEGIN
  NEW.code = (select CONCAT((select value as adjective from internal.adjectives order by random() limit 1), '-', (select value as colour from internal.colours order by random() limit 1), '-', (select value as animal from internal.animals order by random() limit 1)) as code);
  RETURN NEW;
END;$$;
ALTER FUNCTION "internal"."setteamcode"() OWNER TO "postgres";
CREATE TRIGGER "SetTeamCode" BEFORE INSERT ON "public"."games" FOR EACH ROW EXECUTE FUNCTION "internal"."setteamcode"();
