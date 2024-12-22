
CREATE TABLE reservations (
	id serial4 NOT NULL,
	user_id int4 NULL,
	room_id int4 NULL,
	start_time timestamp NOT NULL,
	end_time timestamp NOT NULL,
	booking_date timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	room_price int4 NOT NULL,
	snack_price int4 NULL,
	total_price int4 NOT NULL,
	status varchar(15) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT reservations_pkey PRIMARY KEY (id)
);


ALTER TABLE public.reservations ADD CONSTRAINT reservations_room_id_fkey FOREIGN KEY (room_id) REFERENCES rooms(id);
ALTER TABLE public.reservations ADD CONSTRAINT reservations_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);