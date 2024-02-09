

CREATE TABLE records (
    DoctorId INTEGER NOT NULL,
    PatientId INTEGER, -- was ClientId
    Time_start TIMESTAMP NOT NULL,
    Time_end TIMESTAMP NOT NULL ,


    FOREIGN KEY (DoctorID) REFERENCES doctors (Id),
    FOREIGN KEY (PatientId) REFERENCES patients (Id) -- was FOREIGN KEY (ClientId) REFERENCES clients (Id)
);

INSERT INTO records (DoctorId, PatientId, Time_start, Time_end) VALUES (2, 1, '2024-01-26 10:30', '2024-01-26 10:45');
INSERT INTO records (DoctorId, PatientId, Time_start, Time_end) VALUES (1, 2, '2024-01-26 10:30', '2024-01-26 10:45');

/*
 ALTER TABLE records RENAME COLUMN clientid TO patientid;
 ALTER TABLE records RENAME CONSTRAINT records_clientid_fkey TO records_patientid_fkey;
 */