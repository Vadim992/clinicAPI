CREATE TABLE patients (  -- previous name was 'clients'
    Id SERIAL PRIMARY KEY,
    FirstName VARCHAR NOT NULL, -- was VARCHAR(50)
    LastName VARCHAR NOT NULL, -- was VARCHAR(50)
    Email VARCHAR(320)  UNIQUE , -- was VARCHAR(100), before email was not UNIQUE
    Address VARCHAR NOT NULL -- was VARCHAR(250)
);

INSERT INTO patients (FirstName, LastName, Email, Address) VALUES ('Vadim', 'Pushtakov', 'myemal@gmail.com', 'myadr');
INSERT INTO patients (FirstName, LastName, Email, Address) VALUES ('Ann', 'Lastname', 'heremal@mail.ru', 'heradr');


/*
ALTER TABLE clients RENAME TO patients;
ALTER TABLE patients ADD CONSTRAINT patients_email_key UNIQUE (email);
ALTER TABLE patients RENAME CONSTRAINT clients_pkey TO patients_pkey;
ALTER SEQUENCE clients_id_seq RENAME TO patients_id_seq;

--USE VARCHAR
ALTER TABLE patients ALTER COLUMN firstname TYPE VARCHAR;
ALTER TABLE patients ALTER COLUMN lastname TYPE VARCHAR;
ALTER TABLE patients ALTER COLUMN email TYPE VARCHAR(320);
ALTER TABLE patients ALTER COLUMN address TYPE VARCHAR;
 */


