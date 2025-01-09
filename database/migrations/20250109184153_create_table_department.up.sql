CREATE TABLE IF NOT EXISTS department (
	departmentid varchar(255) NOT NULL DEFAULT gen_random_uuid(),
	departmentname varchar(255) NOT NULL,
	isdeleted bool NOT NULL DEFAULT false,
	createdon timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedon timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	managerid varchar(255) NULL,
	CONSTRAINT department_pkey1 PRIMARY KEY (departmentid),
    CONSTRAINT fk_manager FOREIGN KEY (managerid) REFERENCES manager(managerid)
);
