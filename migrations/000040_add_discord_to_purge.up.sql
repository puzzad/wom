DROP TRIGGER IF EXISTS "announce_game_delete" ON "public"."games";

DROP FUNCTION IF EXISTS "internal"."announce_game_delete";
CREATE OR REPLACE FUNCTION "internal"."announce_game_delete"() RETURNS TRIGGER
    LANGUAGE "plpgsql"
    SECURITY DEFINER AS
$$
BEGIN
    EXECUTE internal.announce(FORMAT('💀 Deleted game: `%s`', OLD.code));
    RETURN NULL;
END
$$;

ALTER FUNCTION "internal"."announce_game_delete"() OWNER TO "postgres";

CREATE OR REPLACE TRIGGER "announce_game_delete"
    AFTER DELETE
    ON "public"."games"
    FOR EACH ROW
EXECUTE FUNCTION "internal"."announce_game_delete"();

DROP TRIGGER IF EXISTS "announce_user_delete" ON "auth"."users";

DROP FUNCTION IF EXISTS "internal"."announce_user_delete";
CREATE OR REPLACE FUNCTION "internal"."announce_user_delete"() RETURNS TRIGGER
    LANGUAGE "plpgsql"
    SECURITY DEFINER AS
$$
BEGIN
    EXECUTE internal.announce(FORMAT('💀 Deleted user: `%s`', OLD.id));
    RETURN NULL;
END
$$;
ALTER FUNCTION "internal"."announce_user_delete"() OWNER TO "postgres";

CREATE OR REPLACE TRIGGER "announce_user_delete"
    AFTER DELETE
    ON "auth"."users"
    FOR EACH ROW
EXECUTE FUNCTION "internal"."announce_user_delete"();