ALTER TABLE pool_type_defs
    ADD COLUMN desc_id varchar;

UPDATE pool_type_defs
SET desc_id = gen_random_uuid()::text;

INSERT INTO pool_desc_roots (desc_id, desc_rn)
SELECT desc_id, 1
FROM pool_type_defs
WHERE desc_id IS NOT NULL
ON CONFLICT (desc_id) DO NOTHING;


ALTER TABLE pool_type_defs
    ALTER COLUMN desc_id SET NOT NULL;

ALTER TABLE pool_type_defs
    ADD CONSTRAINT pool_type_defs_desc_root_fk
    FOREIGN KEY (desc_id) REFERENCES pool_desc_roots (desc_id);

ALTER TABLE proc_type_defs
    ADD COLUMN desc_id varchar;

UPDATE proc_type_defs
SET desc_id = gen_random_uuid()::text;

INSERT INTO proc_desc_roots (desc_id, desc_rn)
SELECT desc_id, 1
FROM proc_type_defs
WHERE desc_id IS NOT NULL
ON CONFLICT (desc_id) DO NOTHING;


ALTER TABLE proc_type_defs
    ALTER COLUMN desc_id SET NOT NULL;

ALTER TABLE proc_type_defs
    ADD CONSTRAINT proc_type_defs_desc_root_fk
    FOREIGN KEY (desc_id) REFERENCES proc_desc_roots (desc_id);

ALTER TABLE pool_type_defs
    DROP COLUMN type_rn;

ALTER TABLE proc_type_defs
    DROP COLUMN type_rn;
