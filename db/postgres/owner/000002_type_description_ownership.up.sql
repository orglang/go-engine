ALTER TABLE pool_type_defs
    ADD COLUMN desc_id varchar NOT NULL;

ALTER TABLE pool_type_defs
    ADD CONSTRAINT pool_type_defs_desc_root_fk
    FOREIGN KEY (desc_id) REFERENCES pool_desc_roots (desc_id);

ALTER TABLE proc_type_defs
    ADD COLUMN desc_id varchar NOT NULL;

ALTER TABLE proc_type_defs
    ADD CONSTRAINT proc_type_defs_desc_root_fk
    FOREIGN KEY (desc_id) REFERENCES proc_desc_roots (desc_id);

ALTER TABLE pool_type_defs
    DROP COLUMN type_rn;

ALTER TABLE proc_type_defs
    DROP COLUMN type_rn;
