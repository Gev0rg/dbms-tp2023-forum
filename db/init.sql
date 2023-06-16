CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE IF NOT EXISTS users
(
    user_id  bigserial,
    nickname citext COLLATE "ucs_basic" NOT NULL UNIQUE PRIMARY KEY,
    fullname text                       NOT NULL,
    about    text,
    email    citext                     NOT NULL UNIQUE
);

CREATE UNLOGGED TABLE IF NOT EXISTS forums
(
    forum_id       bigserial,
    user           citext NOT NULL REFERENCES users (nickname),
    slug           citext NOT NULL PRIMARY KEY,
    title          text   NOT NULL,
    posts          int DEFAULT 0,
    threads        int DEFAULT 0
);

CREATE UNLOGGED TABLE IF NOT EXISTS threads
(
    thread_id bigserial PRIMARY KEY NOT NULL,
    author    citext                NOT NULL REFERENCES users (nickname),
    forum     citext                NOT NULL REFERENCES forums (slug),
    title     text                  NOT NULL,
    message   text                  NOT NULL,
    votes     integer                  DEFAULT 0,
    slug      citext,
    created   timestamp with time zone DEFAULT now()
);

CREATE UNLOGGED TABLE IF NOT EXISTS posts
(
    post_id   bigserial PRIMARY KEY NOT NULL UNIQUE,
    forum     citext REFERENCES forums (slug),
    thread_id integer REFERENCES threads (thread_id),
    author    citext                NOT NULL REFERENCES users (nickname),
    parent    int                      DEFAULT 0,
    message   text                  NOT NULL,
    is_edited bool                     DEFAULT FALSE,
    created   timestamp with time zone DEFAULT now(),
    path      bigint[]                 DEFAULT ARRAY []::INTEGER[]
);

CREATE UNLOGGED TABLE IF NOT EXISTS users_votes
(
    nickname  citext NOT NULL REFERENCES users (nickname),
    thread_id int    NOT NULL REFERENCES threads (thread_id),
    voice     int    NOT NULL
);

CREATE UNLOGGED TABLE IF NOT EXISTS users_forums
(

    nickname citext COLLATE "ucs_basic" NOT NULL REFERENCES users (nickname),
    forum    citext                     NOT NULL REFERENCES forums (slug),
    fullname text,
    about    text,
    email    citext,
    CONSTRAINT user_forum_key unique (nickname, forum)
);
