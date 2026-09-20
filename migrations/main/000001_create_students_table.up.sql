CREATE SCHEMA students;

CREATE TABLE students.accounts
(
    id              UUID                    PRIMARY KEY DEFAULT uuidv7(),
    version         BIGINT      NOT NULL    DEFAULT 1,
    name            varchar(32) NOT NULL    CHECK (name ~ '^[\u0400-\u042F\u0490][\u0430-\u045F\u0491]{1,31}$'),
    surname         varchar(32) NOT NULL    CHECK (surname ~ '^[\u0400-\u042F\u0490][\u0430-\u045F\u0491]{1,31}$'),
    phone_number    varchar(15) NOT NULL    CHECK (phone_number ~ '^\+?[0-9]{9,12}$'),
    telegram_synced boolean     NOT NULL    DEFAULT false,
    telegram_id     BIGINT                  
        CHECK ((telegram_synced = FALSE AND telegram_id IS NULL) OR
            (telegram_synced = TRUE AND telegram_id IS NOT NULL))
);