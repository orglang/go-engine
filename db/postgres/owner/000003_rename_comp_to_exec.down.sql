ALTER TABLE proc_conn_acts RENAME COLUMN exec_rn TO comp_rn;
ALTER TABLE proc_conn_acts RENAME COLUMN exec_id TO comp_id;
ALTER TABLE pool_conn_acts RENAME COLUMN exec_rn TO comp_rn;
ALTER TABLE pool_conn_acts RENAME COLUMN exec_id TO comp_id;

ALTER TABLE proc_exec_conns RENAME COLUMN client_exec_rn TO client_comp_rn;
ALTER TABLE proc_exec_conns RENAME COLUMN client_exec_id TO client_comp_id;
ALTER TABLE proc_exec_conns RENAME COLUMN provider_exec_rn TO provider_comp_rn;
ALTER TABLE proc_exec_conns RENAME COLUMN provider_exec_id TO provider_comp_id;
ALTER TABLE proc_exec_conns RENAME TO proc_comp_conns;
ALTER TABLE pool_exec_conns RENAME COLUMN client_exec_rn TO client_comp_rn;
ALTER TABLE pool_exec_conns RENAME COLUMN client_exec_id TO client_comp_id;
ALTER TABLE pool_exec_conns RENAME COLUMN provider_exec_rn TO provider_comp_rn;
ALTER TABLE pool_exec_conns RENAME COLUMN provider_exec_id TO provider_comp_id;
ALTER TABLE pool_exec_conns RENAME TO pool_comp_conns;

ALTER TABLE proc_exec_vars RENAME COLUMN exec_rn TO comp_rn;
ALTER TABLE proc_exec_vars RENAME COLUMN exec_id TO comp_id;
ALTER TABLE proc_exec_vars RENAME TO proc_comp_vars;
ALTER TABLE pool_exec_vars RENAME COLUMN exec_rn TO comp_rn;
ALTER TABLE pool_exec_vars RENAME COLUMN exec_id TO comp_id;
ALTER TABLE pool_exec_vars RENAME TO pool_comp_vars;

ALTER TABLE proc_exec_evals RENAME COLUMN exec_rn TO comp_rn;
ALTER TABLE proc_exec_evals RENAME COLUMN exec_id TO comp_id;
ALTER TABLE pool_exec_evals RENAME COLUMN exec_rn TO comp_rn;
ALTER TABLE pool_exec_evals RENAME COLUMN exec_id TO comp_id;
