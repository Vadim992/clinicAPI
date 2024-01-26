CREATE TABLE doctors (
    Id             SERIAL PRIMARY KEY,
    FirstName      VARCHAR(50)  NOT NULL,
    LastName       VARCHAR(50)  NOT NULL,
    Specialization VARCHAR(100) NOT NULL,
    Room           INTEGER      NOT NULL
);

INSERT INTO doctors (FirstName, LastName, Specialization, Room) VALUES ('Bob', 'Last','therapist' , 500);
INSERT INTO doctors (FirstName, LastName, Specialization, Room) VALUES ('Jim', 'lastName', 'ophthalmologist', 404);