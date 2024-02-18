-- cd ~/Go/clinicAPI/internal/database/v2.0/v2.01
-- cp changes.sql /var/lib/postgresql/GoApps/clinicAPI
-- sudo chmod -R 747 /var/lib/postgresql/GoApps/clinicAPI
-- sudo -i -u postgres
-- psql -Upostgres -dclinicapi -c 'SET search_path TO clinicapi' -f /var/lib/postgresql/GoApps/clinicAPI/changes.sql

ALTER TABLE patients ADD COLUMN phone_number VARCHAR(11) UNIQUE;

UPDATE patients  SET phone_number = '81111111111' WHERE id = 1;
UPDATE patients  SET phone_number = '82222222222' WHERE id = 2;