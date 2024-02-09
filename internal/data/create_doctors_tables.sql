CREATE TABLE doctors (
    Id             SERIAL PRIMARY KEY,
    FirstName      VARCHAR  NOT NULL, -- was VARCHAR(50)
    LastName       VARCHAR  NOT NULL, -- was VARCHAR(50)
    Specialization VARCHAR NOT NULL, -- was VARCHAR(100)
    Email VARCHAR(320) UNIQUE NOT NULL , -- add this data
    Room           INTEGER      NOT NULL
);

INSERT INTO doctors (FirstName, LastName, Specialization, Email, Room) VALUES ('Bob', 'Last','therapist', 'doc1@mail.ru',500);
INSERT INTO doctors (FirstName, LastName, Specialization,Email, Room) VALUES ('Jim', 'lastName', 'ophthalmologist', 'doc2@mail.ru',404);


/*
ALTER TABLE doctors ADD COLUMN email VARCHAR(320) UNIQUE;


--USE VARCHAR
ALTER TABLE doctors ALTER COLUMN firstname TYPE VARCHAR;
ALTER TABLE doctors ALTER COLUMN lastname TYPE VARCHAR;
ALTER TABLE doctors ALTER COLUMN specialization TYPE VARCHAR;
 */