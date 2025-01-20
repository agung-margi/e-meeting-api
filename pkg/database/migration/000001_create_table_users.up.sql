CREATE TABLE users(
    id SERIAL NOT NULL,
    username varchar(100) NOT NULL,
    email varchar(100) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    is_admin boolean DEFAULT false,
    img_url varchar(255) NULL,
    is_active boolean DEFAULT true,
    language varchar(10) DEFAULT 'english'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);