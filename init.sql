DROP TABLE IF EXISTS person;

CREATE TABLE person (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT, 
    age INT,
    gender TEXT,
    nationality TEXT
);

CREATE INDEX idx_person_name_surname ON person (name, surname);