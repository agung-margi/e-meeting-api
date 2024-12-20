CREATE TABLE reservation_details (
	id serial4 NOT NULL,
	reservation_id int4 NULL,
	"name" varchar(255) NOT NULL,
	phone varchar(15) NOT NULL,
	company varchar(50) NOT NULL,
	snack_id int4 NULL,
	participants int4 NOT NULL,
	notes varchar(255) NULL,
	CONSTRAINT reservation_details_pkey PRIMARY KEY (id)
);


-- public.reservation_details foreign keys

ALTER TABLE public.reservation_details ADD CONSTRAINT reservation_details_resevation_id_fkey FOREIGN KEY (reservation_id) REFERENCES reservations(id) ON DELETE CASCADE;
ALTER TABLE public.reservation_details ADD CONSTRAINT reservation_details_snack_id_fkey FOREIGN KEY (snack_id) REFERENCES snacks(id);