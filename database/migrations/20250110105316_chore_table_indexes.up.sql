
CREATE unique index employees_identitynumber_unique_idx 
ON employees (identitynumber);
CREATE INDEX employees_identitynumber_idx_btree
ON employees USING BTREE (identitynumber);

CREATE INDEX employees_name_idx_fulltext_gin 
ON employees USING GIN (to_tsvector('simple', name));

CREATE INDEX department_name_idx_fulltext_gin 
ON department USING GIN (to_tsvector('simple', departmentname));

CREATE INDEX employees_gender_idx_btree
ON employees USING BTREE (gender);