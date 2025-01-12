CREATE TABLE rooms(
    id SERIAL NOT NULL,
    name varchar(100) NOT NULL,
    room_type_id integer NOT NULL,
    capacity integer NOT NULL,
    price integer NOT NULL,
    img_url varchar(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    CONSTRAINT rooms_room_type_id_fkey FOREIGN key(room_type_id) REFERENCES room_types(id)
);