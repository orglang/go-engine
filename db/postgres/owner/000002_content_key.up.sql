ALTER TABLE pool_type_defs RENAME COLUMN exp_vk TO exp_ck;
ALTER TABLE pool_type_exps RENAME COLUMN exp_vk TO exp_ck;
ALTER TABLE pool_type_exps RENAME COLUMN sup_exp_vk TO sup_exp_ck;
ALTER TABLE pool_comp_vars RENAME COLUMN exp_vk TO exp_ck;

ALTER TABLE proc_type_defs RENAME COLUMN exp_vk TO exp_ck;
ALTER TABLE proc_type_exps RENAME COLUMN exp_vk TO exp_ck;
ALTER TABLE proc_type_exps RENAME COLUMN sup_exp_vk TO sup_exp_ck;
ALTER TABLE proc_comp_vars RENAME COLUMN exp_vk TO exp_ck;
