CREATE TABLE IF NOT EXISTS manager (
	managerid varchar(255) NOT NULL DEFAULT gen_random_uuid(),
	name varchar(255) NULL,
	email varchar(255) NULL,
	password varchar(255) NULL,
	userimageuri varchar(255) NULL,
	companyname varchar(255) NULL,
	companyimageuri varchar(255) NULL,
	isdeleted bool NOT NULL DEFAULT false,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT manager_email_key1 UNIQUE (email),
	CONSTRAINT manager_pkey1 PRIMARY KEY (managerid)
);