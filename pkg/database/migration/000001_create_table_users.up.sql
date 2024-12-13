CREATE TABLE users (
	id serial4 NOT NULL,
	username varchar(100) NOT NULL,
	email varchar(100) NOT NULL,
	"password" varchar(255) NOT NULL,
	is_admin bool DEFAULT false NULL,
	img_url varchar(255) NOT NULL,
	is_active bool DEFAULT true NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);