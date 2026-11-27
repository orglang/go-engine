ALTER TABLE pool_exec_evals RENAME COLUMN comp_id TO exec_id;
ALTER TABLE pool_exec_evals RENAME COLUMN comp_rn TO exec_rn;
ALTER TABLE proc_exec_evals RENAME COLUMN comp_id TO exec_id;
ALTER TABLE proc_exec_evals RENAME COLUMN comp_rn TO exec_rn;

ALTER TABLE pool_comp_vars RENAME TO pool_exec_vars;
ALTER TABLE pool_exec_vars RENAME COLUMN comp_id TO exec_id;
ALTER TABLE pool_exec_vars RENAME COLUMN comp_rn TO exec_rn;
ALTER TABLE proc_comp_vars RENAME TO proc_exec_vars;
ALTER TABLE proc_exec_vars RENAME COLUMN comp_id TO exec_id;
ALTER TABLE proc_exec_vars RENAME COLUMN comp_rn TO exec_rn;

ALTER TABLE pool_comp_conns RENAME TO pool_exec_conns;
ALTER TABLE pool_exec_conns RENAME COLUMN provider_comp_id TO provider_exec_id;
ALTER TABLE pool_exec_conns RENAME COLUMN provider_comp_rn TO provider_exec_rn;
ALTER TABLE pool_exec_conns RENAME COLUMN client_comp_id TO client_exec_id;
ALTER TABLE pool_exec_conns RENAME COLUMN client_comp_rn TO client_exec_rn;
ALTER TABLE proc_comp_conns RENAME TO proc_exec_conns;
ALTER TABLE proc_exec_conns RENAME COLUMN provider_comp_id TO provider_exec_id;
ALTER TABLE proc_exec_conns RENAME COLUMN provider_comp_rn TO provider_exec_rn;
ALTER TABLE proc_exec_conns RENAME COLUMN client_comp_id TO client_exec_id;
ALTER TABLE proc_exec_conns RENAME COLUMN client_comp_rn TO client_exec_rn;

ALTER TABLE pool_conn_acts RENAME COLUMN comp_id TO exec_id;
ALTER TABLE pool_conn_acts RENAME COLUMN comp_rn TO exec_rn;
ALTER TABLE proc_conn_acts RENAME COLUMN comp_id TO exec_id;
ALTER TABLE proc_conn_acts RENAME COLUMN comp_rn TO exec_rn;
