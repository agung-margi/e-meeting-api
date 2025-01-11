CREATE TABLE reservations(
    id SERIAL NOT NULL,
    user_id integer,
    room_id integer,
    start_time timestamp without time zone NOT NULL,
    end_time timestamp without time zone NOT NULL,
    booking_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    room_price integer NOT NULL,
    snack_price integer,
    total_price integer NOT NULL,
    status varchar(15) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    expired_at timestamp without time zone DEFAULT (CURRENT_TIMESTAMP + INTERVAL '1 hour'),
    PRIMARY KEY(id),
    CONSTRAINT reservations_room_id_fkey FOREIGN key(room_id) REFERENCES rooms(id),
    CONSTRAINT reservations_user_id_fkey FOREIGN key(user_id) REFERENCES users(id)
);