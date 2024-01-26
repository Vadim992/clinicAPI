CREATE TABLE clients (
    Id SERIAL PRIMARY KEY,
    FirstName VARCHAR(50) NOT NULL,
    LastName VARCHAR(50) NOT NULL,
    Email VARCHAR(100),
    Address VARCHAR(250) NOT NULL
);

INSERT INTO clients (FirstName, LastName, Email, Address) VALUES ('Vadim', 'Pushtakov', 'myemal', 'myadr');
INSERT INTO clients (FirstName, LastName, Email, Address) VALUES ('Ann', 'Lastname', 'heremal', 'heradr');