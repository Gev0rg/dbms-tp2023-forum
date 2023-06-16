CREATE UNLOGGED TABLE IF NOT EXISTS users (
    id       bigserial,
    nickname text       NOT NULL UNIQUE PRIMARY KEY,
    fullname text       NOT NULL,
    about    text,
    email    text       NOT NULL UNIQUE
);
