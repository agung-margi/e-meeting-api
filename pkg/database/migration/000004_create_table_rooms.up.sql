CREATE TABLE rooms (
	id serial4 NOT NULL,
	"name" varchar(100) NOT NULL,
	room_type_id int4 NOT NULL,
	capacity int4 NOT NULL,
	price int4 NOT NULL,
	img_url varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT rooms_pkey PRIMARY KEY (id)
);


-- public.rooms foreign keys

ALTER TABLE public.rooms ADD CONSTRAINT rooms_room_type_id_fkey FOREIGN KEY (room_type_id) REFERENCES room_types(id);