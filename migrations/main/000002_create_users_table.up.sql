CREATE SCHEMA users;

CREATE TYPE users.role AS ENUM ('teacher', 'admin');

CREATE TABLE users.accounts
(
    id            UUID                  PRIMARY KEY DEFAULT uuidv7(),
    version       BIGINT       NOT NULL DEFAULT 1,
    role          users.role   NOT NULL DEFAULT 'teacher',
    username      varchar(32)  NOT NULL UNIQUE CHECK (username ~ '^[a-z][a-z0-9._-]{1,30}[a-z0-9]$'),
    name          varchar(32)  NOT NULL CHECK (name ~ '^[\u0400-\u042F\u0490][\u0430-\u045F\u0491]{1,31}$'),
    surname       varchar(32)  NOT NULL CHECK (surname ~ '^[\u0400-\u042F\u0490][\u0430-\u045F\u0491]{1,31}$'),
    password_hash varchar(128) NOT NULL CHECK (password_hash ~ '^\$2[aby]\$(0[4-9]|[12][0-9]|3[01])\$[./A-Za-z0-9]{53}$')
);

ALTER TABLE students.accounts
    ADD COLUMN owner_user_id UUID NOT NULL 
        REFERENCES users.accounts(id) ON DELETE CASCADE;