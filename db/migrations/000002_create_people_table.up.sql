CREATE TABLE IF NOT EXISTS person (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name varchar(100) NOT NULL,
    surname varchar(100) NOT NULL,
    patronymic varchar(100),
    gender varchar(10),
    nationality varchar(2),
    age int,
    CONSTRAINT person_pk PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS person_name_idx ON person (name ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS person_surname_idx ON person (surname ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS person_nationality_idx ON person (nationality ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS person_gender_idx ON person (gender ASC NULLS LAST);
