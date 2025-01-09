-- public.employees definition

-- Drop table

-- DROP TABLE public.employees;

CREATE TABLE public.employees (
	id varchar(255) NOT NULL DEFAULT gen_random_uuid(),
	"name" varchar(255) NOT NULL,
	identitynumber varchar(255) NOT NULL UNIQUE,
	employeeimageuri varchar(255) NOT NULL,
	gender varchar(6) NOT NULL,
	departmentid varchar NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT employees_pkey PRIMARY KEY (id)
);


-- public.employees foreign keys

ALTER TABLE public.employees ADD CONSTRAINT fk_manager FOREIGN KEY (departmentid) REFERENCES public.department(departmentid);

-- Previously created this constraint, but need to remove due to requirement adjustment
-- ALTER TABLE public.employees ADD CONSTRAINT unique_employee_identitynumber UNIQUE (identitynumber);

-- Use this query to remove the unique constarint on identityNumber
-- ALTER TABLE public.employees DROP CONSTRAINT unique_employee_identitynumber;

ALTER TABLE public.employees ALTER COLUMN id SET DEFAULT gen_random_uuid();