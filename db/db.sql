CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;

CREATE TABLE IF NOT EXISTS users (
    nickname    CITEXT COLLATE "C"  NOT NULL PRIMARY KEY,
    fullname    TEXT                NOT NULL,
    email       CITEXT              NOT NULL UNIQUE,
    about       TEXT                NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS forum (
    slug        CITEXT       NOT NULL PRIMARY KEY,
    title       TEXT         NOT NULL,
    author      CITEXT       NOT NULL,
    posts       INTEGER      NOT NULL DEFAULT 0,
    threads     INTEGER      NOT NULL DEFAULT 0,
	FOREIGN KEY (author) REFERENCES users (nickname)
);

CREATE TABLE IF NOT EXISTS threads (
    id          SERIAL                      NOT NULL PRIMARY KEY,
    slug        CITEXT                      DEFAULT '',
    title       TEXT                        NOT NULL,
    author      CITEXT                      NOT NULL,
    forum       CITEXT                      NOT NULL,
    message     TEXT                        NOT NULL,
    votes       INT DEFAULT 0               NOT NULL,
    created     TIMESTAMP WITH TIME ZONE    DEFAULT now(),
    FOREIGN KEY (author) REFERENCES users (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug)
);

CREATE TABLE IF NOT EXISTS posts (
    id          BIGSERIAL                   NOT NULL PRIMARY KEY,
    parent      BIGINT                      NOT NULL,
    author      CITEXT                      NOT NULL,
    message     TEXT                        NOT NULL,
    isEdited    BOOLEAN                     NOT NULL DEFAULT FALSE,
    forum       CITEXT                      NOT NULL,
    thread      INTEGER                     NOT NULL,
    created     TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    path        BIGINT                      ARRAY,
    FOREIGN KEY (author) REFERENCES users (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug),
    FOREIGN KEY (thread) REFERENCES threads (id)
);

CREATE TABLE IF NOT EXISTS votes (
	nickname 	CITEXT	NOT NULL,
  	thread 		INT		NOT NULL,
  	voice     	INT		NOT NULL,
	FOREIGN KEY (nickname) REFERENCES users(nickname),
	FOREIGN KEY (thread) REFERENCES threads(id),
    PRIMARY KEY (nickname, thread)
);

CREATE TABLE IF NOT EXISTS forum_users (
    nickname    CITEXT COLLATE "C"  NOT NULL,
    fullname    TEXT                NOT NULL,
    email       CITEXT              NOT NULL,
    about       TEXT                NOT NULL DEFAULT '',
    forum       CITEXT              NOT NULL,
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug),
	PRIMARY KEY (nickname, forum)
);

CREATE UNIQUE INDEX IF NOT EXISTS index_threads_slug ON threads(slug) WHERE TRIM(slug) <> '';
CREATE INDEX IF NOT EXISTS index_thread__forum_created ON threads(forum, created);
CREATE INDEX IF NOT EXISTS index_thread__slug_id_forum ON threads(slug, id, forum);
CREATE INDEX IF NOT EXISTS index_posts__thread ON posts(thread);
CREATE INDEX IF NOT EXISTS index_posts__id_thread ON posts(id, thread);

CREATE OR REPLACE FUNCTION vote_insert() RETURNS TRIGGER AS $vote_insert$
BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice
    WHERE id = NEW.thread;
    RETURN NULL;
END;
$vote_insert$  LANGUAGE plpgsql;

CREATE TRIGGER vote_insert AFTER INSERT ON votes FOR EACH ROW EXECUTE PROCEDURE vote_insert();

CREATE OR REPLACE FUNCTION vote_update() RETURNS TRIGGER AS $vote_update$
BEGIN
	IF OLD.voice = NEW.voice
		THEN RETURN NULL;
	END IF;
  	UPDATE threads
	SET
		votes = votes + CASE WHEN NEW.voice = -1 THEN -2 ELSE 2 END
  	WHERE id = NEW.thread;
  	RETURN NULL;
END;
$vote_update$ LANGUAGE  plpgsql;

CREATE TRIGGER vote_update AFTER UPDATE ON votes FOR EACH ROW EXECUTE PROCEDURE vote_update();

CREATE OR REPLACE FUNCTION increment_posts_count() RETURNS TRIGGER AS $increment_posts_count$
BEGIN
    UPDATE forum SET 
        posts = (posts + 1)
    WHERE slug = NEW.forum;
    
    RETURN NULL;
END;
$increment_posts_count$ LANGUAGE plpgsql;

CREATE TRIGGER increment_posts_count AFTER INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE increment_posts_count();

CREATE OR REPLACE FUNCTION increment_threads_count() RETURNS TRIGGER AS $increment_threads_count$
BEGIN
    UPDATE forum SET 
        threads = (threads + 1)
    WHERE slug = NEW.forum;
    
    RETURN NULL;
END;
$increment_threads_count$ LANGUAGE plpgsql;

CREATE TRIGGER increment_threads_count AFTER INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE increment_threads_count();

CREATE OR REPLACE FUNCTION post_paste_forum_user() RETURNS TRIGGER AS $post_paste_forum_user$
BEGIN
    INSERT INTO forum_users
    SELECT nickname, fullname, email, about, NEW.forum as forum 
    FROM users
    WHERE nickname = NEW.author
	ON CONFLICT DO NOTHING;
    
    RETURN NULL;
END;
$post_paste_forum_user$ LANGUAGE plpgsql;

CREATE TRIGGER post_paste_forum_user AFTER INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE post_paste_forum_user();

CREATE OR REPLACE FUNCTION thread_paste_forum_user() RETURNS TRIGGER AS $thread_paste_forum_user$
BEGIN
    INSERT INTO forum_users
    SELECT nickname, fullname, email, about, NEW.forum as forum 
    FROM users
    WHERE nickname = NEW.author
	ON CONFLICT DO NOTHING;
    
    RETURN NULL;
END;
$thread_paste_forum_user$ LANGUAGE plpgsql;

CREATE TRIGGER thread_paste_forum_user AFTER INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE thread_paste_forum_user();
