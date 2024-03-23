CREATE TABLE role_permissions (
  role_id INTEGER NOT NULL,
  paths text UNIQUE NOT NULL,
  actions text[] NOT NULL
);

INSERT INTO role_permissions (role_id, paths, actions) VALUES (3, '/patients/id', '{"GET", "PUT", "PATCH", "DELETE"}');
-- INSERT INTO role_permissions (role_id, paths, actions) VALUES (3, '/patients', '{"POST"}');