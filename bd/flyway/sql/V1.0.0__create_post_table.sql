CREATE TABLE IF NOT EXISTS Post
(
    id        SERIAL,
    user_id   INT       NOT NULL,
    content   TEXT      NOT NULL,
    is_open   BOOLEAN   NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    PRIMARY KEY (id, timestamp)
) PARTITION BY RANGE (timestamp);