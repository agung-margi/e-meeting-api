CREATE TABLE reset_token (
    id integer GENERATED ALWAYS AS IDENTITY NOT NULL,
    email varchar(100) NOT NULL,
    token varchar(255) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at timestamp without time zone,
    PRIMARY KEY (id),
    CONSTRAINT fk_email FOREIGN KEY (email) REFERENCES users (email) ON DELETE CASCADE
);