CREATE TABLE IF NOT EXISTS Comment
(
    id               SERIAL PRIMARY KEY ,
    parent_id        INT       NULL,
    post_id          INT       NOT NULL,
    user_id          INT       NOT NULL,
    content          TEXT      NOT NULL,
    timestamp        TIMESTAMP NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES Comment (id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES Post (id) ON DELETE CASCADE
);