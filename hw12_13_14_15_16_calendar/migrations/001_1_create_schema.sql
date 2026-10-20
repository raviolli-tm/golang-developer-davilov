CREATE SCHEMA if not exists calendar AUTHORIZATION postgres_user;

GRANT USAGE ON SCHEMA calendar TO postgres_user;
GRANT USAGE ON SCHEMA calendar TO notify_user;

GRANT CREATE ON SCHEMA calendar TO postgres_user;