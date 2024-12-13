CREATE TABLE snacks (
	id serial4 NOT NULL,
	"name" varchar(100) NOT NULL,
	price int4 NOT NULL,
	CONSTRAINT snacks_pkey PRIMARY KEY (id)
);