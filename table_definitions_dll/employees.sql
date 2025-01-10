-- public.employees definition

-- Drop table

-- DROP TABLE public.employees;

CREATE TABLE public.employees (
	id varchar(255) NOT NULL DEFAULT gen_random_uuid(),
	"name" varchar(255) NOT NULL,
	identitynumber varchar(255) NOT NULL,
	employeeimageuri varchar(255) NOT NULL,
	gender varchar(6) NOT NULL,
	departmentid varchar NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT employees_identitynumber_key UNIQUE (identitynumber),
	CONSTRAINT employees_pkey PRIMARY KEY (id)
);
CREATE INDEX employees_gender_idx_btree ON public.employees USING btree (gender);
CREATE INDEX employees_identitynumber_idx_btree ON public.employees USING btree (identitynumber);
CREATE INDEX employees_name_idx_fulltext_gin ON public.employees USING gin (to_tsvector('simple'::regconfig, (name)::text));

-- public.employees foreign keys
ALTER TABLE public.employees ADD CONSTRAINT fk_manager1 FOREIGN KEY (departmentid) REFERENCES public.department(departmentid);