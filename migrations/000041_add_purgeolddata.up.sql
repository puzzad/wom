ALTER TABLE games
    DROP CONSTRAINT IF EXISTS games_user_fkey,
    ADD CONSTRAINT games_user_fkey FOREIGN KEY ("user") REFERENCES auth.users ON DELETE CASCADE;
ALTER TABLE guesses
    DROP CONSTRAINT IF EXISTS guesses_game_fkey,
    ADD CONSTRAINT guesses_game_fkey FOREIGN KEY (game) REFERENCES games (code) ON DELETE CASCADE;
ALTER TABLE internal.usedhints
    DROP CONSTRAINT IF EXISTS usedhints_game_fkey,
    ADD CONSTRAINT usedhints_game_fkey FOREIGN KEY (game) REFERENCES games (code) ON DELETE CASCADE;
DROP FUNCTION IF EXISTS "internal"."purgeolddata";
CREATE OR REPLACE FUNCTION "internal"."purgeolddata"()
    RETURNS TABLE
            (
                "users"        INT,
                "games"        INT,
                "transactions" INT
            )
    LANGUAGE "plpgsql"
    SECURITY DEFINER
AS
$$
DECLARE
    userPurgeDays        INT := 365;
    gamePurgeDays        INT := 365;
    transactionPurgeDays INT := 2555;
    data                 INT;
BEGIN
    SELECT "value" FROM "internal"."secrets" WHERE name = 'userPurgeDays' INTO data;
    IF data IS NOT NULL THEN
        SELECT data::INT INTO userPurgeDays;
    END IF;
    SELECT "value" FROM "internal"."secrets" WHERE name = 'gamePurgeDays' INTO data;
    IF data IS NOT NULL THEN
        SELECT data::INT INTO gamePurgeDays;
    END IF;
    SELECT "value" FROM "internal"."secrets" WHERE name = 'transactionPurgeDays' INTO data;
    IF data IS NOT NULL THEN
        SELECT data::INT INTO transactionPurgeDays;
    END IF;
    DELETE FROM "auth"."users" WHERE last_sign_in_at < NOW()::DATE - userPurgeDays;
    DELETE FROM "public"."games" WHERE "endTime" < NOW()::DATE - gamePurgeDays;
END
$$;
ALTER FUNCTION "internal"."purgeolddata"() OWNER TO "postgres";
SELECT cron.schedule('purgeolddata', '0 3 * * *', 'SELECT "internal"."purgeolddata"()');
