CREATE TABLE reservation_details(
    id SERIAL NOT NULL,
    reservation_id integer,
    name varchar(255) NOT NULL,
    phone varchar(15) NOT NULL,
    company varchar(50) NOT NULL,
    snack_id integer,
    participants integer NOT NULL,
    notes varchar(255),
    PRIMARY KEY(id),
    CONSTRAINT reservation_details_resevation_id_fkey FOREIGN key(reservation_id) REFERENCES reservations(id),
    CONSTRAINT reservation_details_snack_id_fkey FOREIGN key(snack_id) REFERENCES snacks(id)
);