CREATE TABLE IF NOT EXISTS Comment
(
    id        SERIAL,
    parent_id INT,
    post_id   INT           NOT NULL,
    user_id   INT           NOT NULL,
    content   VARCHAR(2000) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    PRIMARY KEY (id, timestamp),
    CONSTRAINT fk_post
        FOREIGN KEY (post_id, timestamp)
            REFERENCES Post (id, timestamp)
            ON DELETE CASCADE,
    CONSTRAINT fk_parent
        FOREIGN KEY (parent_id, timestamp)
            REFERENCES Comment (id, timestamp)
            ON DELETE CASCADE
) PARTITION BY RANGE (timestamp);