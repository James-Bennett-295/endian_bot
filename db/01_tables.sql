-- Create tables

CREATE TABLE users (
    id              SERIAL          PRIMARY KEY,
    discord_id      VARCHAR(20)     UNIQUE          NOT NULL
);

CREATE TABLE timezones (
    id              SMALLSERIAL     PRIMARY KEY,
    name            VARCHAR(32)     UNIQUE          NOT NULL
);

CREATE TABLE user_timezones (
    user_id         INTEGER         UNIQUE          NOT NULL,
    timezone_id     SMALLINT                        NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE RESTRICT,
    FOREIGN KEY (timezone_id)
        REFERENCES timezones (id)
        ON DELETE RESTRICT
);
