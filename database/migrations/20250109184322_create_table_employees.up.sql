CREATE TABLE IF NOT EXISTS employees (
	id varchar(255) NOT NULL DEFAULT gen_random_uuid(),
	name varchar(255) NOT NULL,
	identitynumber varchar(255) NOT NULL UNIQUE,
	employeeimageuri varchar(255) NOT NULL,
	gender varchar(6) NOT NULL,
	departmentid varchar NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT employees_pkey PRIMARY KEY (id),
    CONSTRAINT fk_manager1 FOREIGN KEY (departmentid) REFERENCES department (departmentid)
);