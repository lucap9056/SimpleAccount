CREATE TABLE User (
    id          INT             AUTO_INCREMENT,
    username    VARCHAR(32)     NOT NULL,
    email       VARCHAR(256)    NOT NULL,
    salt        VARCHAR(64)     NOT NULL,
    hash        VARCHAR(64)     NOT NULL,
    last_edit   TIMESTAMP       NOT NULL        DEFAULT CURRENT_TIMESTAMP(),
    create_time TIMESTAMP       NOT NULL        DEFAULT CURRENT_TIMESTAMP(),
    deleted     TIMESTAMP,
    PRIMARY KEY (id)
);