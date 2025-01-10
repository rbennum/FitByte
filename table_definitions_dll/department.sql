-- public.department definition

-- Drop table

-- DROP TABLE public.department;

CREATE TABLE public.department (
	departmentid varchar(255) NOT NULL DEFAULT gen_random_uuid(),
	departmentname varchar(255) NOT NULL,
	isdeleted bool NOT NULL DEFAULT false,
	createdon timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedon timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	managerid varchar(255) NULL,
	CONSTRAINT department_pkey1 PRIMARY KEY (departmentid)
);
CREATE INDEX department_name_idx_fulltext_gin ON public.department USING gin (to_tsvector('simple'::regconfig, (departmentname)::text));

-- public.department foreign keys
ALTER TABLE public.department ADD CONSTRAINT fk_manager FOREIGN KEY (managerid) REFERENCES public.manager(managerid);