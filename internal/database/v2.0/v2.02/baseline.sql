--
-- PostgreSQL database dump
--

-- Dumped from database version 16.1 (Ubuntu 16.1-1.pgdg22.04+1)
-- Dumped by pg_dump version 16.1 (Ubuntu 16.1-1.pgdg22.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: clinicapi; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA clinicapi;


ALTER SCHEMA clinicapi OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: doctors; Type: TABLE; Schema: clinicapi; Owner: postgres
--

CREATE TABLE clinicapi.doctors (
    id integer NOT NULL,
    firstname character varying NOT NULL,
    lastname character varying NOT NULL,
    specialization character varying NOT NULL,
    room integer NOT NULL,
    email character varying(320) NOT NULL
);


ALTER TABLE clinicapi.doctors OWNER TO postgres;

--
-- Name: doctors_id_seq; Type: SEQUENCE; Schema: clinicapi; Owner: postgres
--

CREATE SEQUENCE clinicapi.doctors_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE clinicapi.doctors_id_seq OWNER TO postgres;

--
-- Name: doctors_id_seq; Type: SEQUENCE OWNED BY; Schema: clinicapi; Owner: postgres
--

ALTER SEQUENCE clinicapi.doctors_id_seq OWNED BY clinicapi.doctors.id;


--
-- Name: patients; Type: TABLE; Schema: clinicapi; Owner: postgres
--

CREATE TABLE clinicapi.patients (
    id integer NOT NULL,
    firstname character varying NOT NULL,
    lastname character varying NOT NULL,
    email character varying(320),
    address character varying NOT NULL,
    phone_number character varying(11)
);


ALTER TABLE clinicapi.patients OWNER TO postgres;

--
-- Name: patients_id_seq; Type: SEQUENCE; Schema: clinicapi; Owner: postgres
--

CREATE SEQUENCE clinicapi.patients_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE clinicapi.patients_id_seq OWNER TO postgres;

--
-- Name: patients_id_seq; Type: SEQUENCE OWNED BY; Schema: clinicapi; Owner: postgres
--

ALTER SEQUENCE clinicapi.patients_id_seq OWNED BY clinicapi.patients.id;


--
-- Name: records; Type: TABLE; Schema: clinicapi; Owner: postgres
--

CREATE TABLE clinicapi.records (
    doctorid integer NOT NULL,
    patientid integer,
    time_start timestamp without time zone NOT NULL,
    time_end timestamp without time zone NOT NULL
);


ALTER TABLE clinicapi.records OWNER TO postgres;

--
-- Name: doctors id; Type: DEFAULT; Schema: clinicapi; Owner: postgres
--

ALTER TABLE ONLY clinicapi.doctors ALTER COLUMN id SET DEFAULT nextval('clinicapi.doctors_id_seq'::regclass);


--
-- Name: patients id; Type: DEFAULT; Schema: clinicapi; Owner: postgres
--

ALTER TABLE ONLY clinicapi.patients ALTER COLUMN id SET DEFAULT nextval('clinicapi.patients_id_seq'::regclass);


--
-- Data for Name: doctors; Type: TABLE DATA; Schema: clinicapi; Owner: postgres
--

COPY clinicapi.doctors (id, firstname, lastname, specialization, room, email) FROM stdin;
1	Bob	Last	therapist	500	doc1@mail.ru
4	Ann	Doctorova	ophthalmologist	201	new1@gmail.com
8	Diana	Doctorova	surgeon	330	sd1@gmail.com
5	Diana	Doctorova	surgeon	202	sd@gmail.com
2	Jim	lastName	therapist	404	doc2@mail.ru
\.


--
-- Data for Name: patients; Type: TABLE DATA; Schema: clinicapi; Owner: postgres
--

COPY clinicapi.patients (id, firstname, lastname, email, address, phone_number) FROM stdin;
1	Vadim	Pushtakov	myemal@gmail.com	myadr	81111111111
2	Ann	Lastname	heremal@mail.ru	heradr	82222222222
5	Andrey	Push	him@mail.ru	addr	81231231214
8	Name	Push	\N	addr	81231231213
\.


--
-- Data for Name: records; Type: TABLE DATA; Schema: clinicapi; Owner: postgres
--

COPY clinicapi.records (doctorid, patientid, time_start, time_end) FROM stdin;
1	\N	2024-04-01 09:00:00	2024-04-01 09:15:00
2	\N	2024-04-01 09:00:00	2024-04-01 09:15:00
4	\N	2024-04-01 09:00:00	2024-04-01 09:15:00
5	\N	2024-04-01 09:00:00	2024-04-01 09:15:00
8	\N	2024-04-01 09:00:00	2024-04-01 09:15:00
1	1	2024-04-01 09:15:00	2024-04-01 09:30:00
1	2	2024-04-01 09:30:00	2024-04-01 09:45:00
4	2	2024-04-01 09:45:00	2024-04-01 10:00:00
8	5	2024-04-01 09:45:00	2024-04-01 10:00:00
5	8	2024-04-01 09:45:00	2024-04-01 10:00:00
2	\N	2024-04-01 10:00:00	2024-04-01 10:15:00
\.


--
-- Name: doctors_id_seq; Type: SEQUENCE SET; Schema: clinicapi; Owner: postgres
--

SELECT pg_catalog.setval('clinicapi.doctors_id_seq', 8, true);


--
-- Name: patients_id_seq; Type: SEQUENCE SET; Schema: clinicapi; Owner: postgres
--

SELECT pg_catalog.setval('clinicapi.patients_id_seq', 8, true);


--
-- Name: doctors doctors_email_key; Type: CONSTRAINT; Schema: clinicapi; Owner: postgres
--

ALTER TABLE ONLY clinicapi.doctors
    ADD CONSTRAINT doctors_email_key UNIQUE (email);


--
-- Name: doctors doctors_pkey; Type: CONSTRAINT; Schema: clinicapi; Owner: postgres
--

ALTER TABLE ONLY clinicapi.doctors
    ADD CONSTRAINT doctors_pkey PRIMARY KEY (id);


--
-- Name: patients patients_email_key; Type: CONSTRAINT; Schema: clinicapi; Owner: postgres
--

ALTER TABLE ONLY clinicapi.patients
    ADD CONSTRAINT patients_email_key UNIQUE (email);


--
-- Name: patients patients_phone_number_key; Type: CONSTRAINT; Schema: clinicapi; Owner: postgres
--

ALTER TABLE ONLY clinicapi.patients
    ADD CONSTRAINT patients_phone_number_key UNIQUE (phone_number);


--
-- Name: patients patients_pkey; Type: CONSTRAINT; Schema: clinicapi; Owner: postgres
--

ALTER TABLE ONLY clinicapi.patients
    ADD CONSTRAINT patients_pkey PRIMARY KEY (id);


--
-- Name: records records_doctorid_fkey; Type: FK CONSTRAINT; Schema: clinicapi; Owner: postgres
--

ALTER TABLE ONLY clinicapi.records
    ADD CONSTRAINT records_doctorid_fkey FOREIGN KEY (doctorid) REFERENCES clinicapi.doctors(id);


--
-- Name: records records_patientid_fkey; Type: FK CONSTRAINT; Schema: clinicapi; Owner: postgres
--

ALTER TABLE ONLY clinicapi.records
    ADD CONSTRAINT records_patientid_fkey FOREIGN KEY (patientid) REFERENCES clinicapi.patients(id);


--
-- Name: SCHEMA clinicapi; Type: ACL; Schema: -; Owner: postgres
--

GRANT USAGE ON SCHEMA clinicapi TO clinicapi_user;


--
-- Name: TABLE doctors; Type: ACL; Schema: clinicapi; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE clinicapi.doctors TO clinicapi_user;


--
-- Name: SEQUENCE doctors_id_seq; Type: ACL; Schema: clinicapi; Owner: postgres
--

GRANT SELECT,USAGE ON SEQUENCE clinicapi.doctors_id_seq TO clinicapi_user;


--
-- Name: TABLE patients; Type: ACL; Schema: clinicapi; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE clinicapi.patients TO clinicapi_user;


--
-- Name: SEQUENCE patients_id_seq; Type: ACL; Schema: clinicapi; Owner: postgres
--

GRANT SELECT,USAGE ON SEQUENCE clinicapi.patients_id_seq TO clinicapi_user;


--
-- Name: TABLE records; Type: ACL; Schema: clinicapi; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE clinicapi.records TO clinicapi_user;


--
-- PostgreSQL database dump complete
--

