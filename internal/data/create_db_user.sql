CREATE DATABASE clinicapi WITH ENCODING='UTF8' LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8' TEMPLATE=template0;

CREATE USER clinicapi_user PASSWORD 'mypass';

-- \c clinicapi
CREATE SCHEMA clinicapi;
GRANT USAGE ON SCHEMA clinicapi to clinicapi_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA clinicapi TO clinicapi_user;
GRANT SELECT, UPDATE, INSERT, DELETE ON ALL TABLES IN SCHEMA clinicapi to clinicapi_user;
SET search_path TO clinicapi;