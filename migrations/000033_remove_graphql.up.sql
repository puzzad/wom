DROP EXTENSION IF EXISTS pg_graphql;
DROP EVENT TRIGGER IF EXISTS graphql_watch_ddl;
DROP EVENT TRIGGER IF EXISTS graphql_watch_drop;
DROP EVENT TRIGGER IF EXISTS issue_graphql_placeholder;
DROP EVENT TRIGGER IF EXISTS issue_pg_graphql_access;
DROP SCHEMA IF EXISTS graphql_public CASCADE;