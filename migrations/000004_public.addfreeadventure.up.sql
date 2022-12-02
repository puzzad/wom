CREATE FUNCTION "public"."addfreeadventure"("adventureid" bigint) RETURNS character varying
    LANGUAGE "plpgsql" SECURITY DEFINER
    AS $$DECLARE
 res varchar;
BEGIN
if auth.uid() is null then
raise exception 'not logged in';
end if;
if not exists(select 1 from adventures where id = adventureId and price = 0) then
raise exception 'free adventure not found';
end if;
insert into games(status, "user", adventure) values ('PAID', auth.uid(), adventureId) returning code into res;
return res;
END$$;
ALTER FUNCTION "public"."addfreeadventure"("adventureid" bigint) OWNER TO "postgres";
