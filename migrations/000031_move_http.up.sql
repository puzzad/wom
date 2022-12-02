DO $$
DECLARE
  exists integer := 0;
BEGIN
  SELECT count(*) FROM pg_extension WHERE extname='pg_http' into exists;
  IF exists > 0 THEN
    ALTER EXTENSION http SET SCHEMA extensions;
  END IF;
END $$