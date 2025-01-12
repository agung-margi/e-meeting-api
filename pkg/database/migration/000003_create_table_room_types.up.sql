CREATE TABLE room_types(
    id SERIAL NOT NULL,
    name varchar(100) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);