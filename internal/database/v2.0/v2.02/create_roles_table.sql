CREATE TABLE roles (
    id        SERIAL PRIMARY KEY,
    role_name VARCHAR NOT NULL
);

INSERT INTO roles (role_name) VALUES ('admin');
INSERT INTO roles (role_name) VALUES ('doctor');
INSERT INTO roles (role_name) VALUES ('patient');