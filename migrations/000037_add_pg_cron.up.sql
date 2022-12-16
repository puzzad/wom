CREATE EXTENSION IF NOT EXISTS pg_cron;
SELECT cron.schedule('nightly-vacuum', '0 2 * * *', 'VACUUM');
SELECT cron.schedule('purgeolddata', '0 3 * * *', 'SELECT "internal"."purgeolddata"()');
