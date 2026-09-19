CREATE SCHEMA students;

CREATE TABLE students.accounts 
(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    version BIGINT NOT NULL DEFAULT 1,
    
    name varchar(20) NOT NULL,
    surname varchar(20) NOT NULL,
    phone_number varchar(20) NOT NULL,
    telegram_id varchar(50)
);