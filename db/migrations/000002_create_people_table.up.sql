CREATE TABLE IF NOT EXISTS people (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name varchar(100) NOT NULL,
    surname varchar(100) NOT NULL,
    patronymic varchar(100),
    age int,
    gender varchar(10),
    nationality varchar(2),
    is_deleted boolean DEFAULT false NOT NULL,
    CONSTRAINT people_pk PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS people_name_idx ON people (name ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS people_surname_idx ON people (surname ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS people_nationality_idx ON people (nationality ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS people_gender_idx ON people (gender ASC NULLS LAST);
