-- +goose Up

CREATE TABLE IF NOT EXISTS userprofile (
	user_id varchar(20) NOT NULL,
	password varchar(100) NOT NULL,
	first_name varchar(150) NOT NULL,
	last_name varchar(150) NOT NULL,
	address text NULL,
	birthdate date NULL,
	email varchar(255) NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT user_id_pkey PRIMARY KEY (user_id)
);

-- +goose Down
DROP TABLE userprofile;