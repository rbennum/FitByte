-- public.employees definition

-- Drop table

-- DROP TABLE public.employees;

CREATE TABLE public.employees (
	id varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	identitynumber varchar(255) NOT NULL,
	employeeimageuri varchar(255) NOT NULL,
	gender varchar(6) NOT NULL,
	departmentid varchar NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT employees_pkey PRIMARY KEY (id)
);


-- public.employees foreign keys

ALTER TABLE public.employees ADD CONSTRAINT fk_manager FOREIGN KEY (departmentid) REFERENCES public.department(departmentid);