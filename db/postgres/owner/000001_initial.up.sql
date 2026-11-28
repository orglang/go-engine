CREATE EXTENSION IF NOT EXISTS ltree;

CREATE TABLE pool_desc_roots (
    desc_id varchar PRIMARY KEY,
    desc_rn bigint NOT NULL
);

CREATE TABLE pool_desc_binds (
    desc_qn ltree UNIQUE,
    desc_id varchar,
    kind smallint
);

CREATE INDEX pool_desc_qn_gist_idx ON pool_desc_binds USING GIST (desc_qn);

CREATE TABLE pool_type_defs (
    type_id varchar UNIQUE,
    desc_id varchar NOT NULL,
    exp_ck bigint
);

CREATE TABLE pool_type_exps (
    exp_ck bigint UNIQUE,
    sup_exp_ck bigint,
    kind smallint,
    spec jsonb
);

CREATE TABLE pool_desc_terms (
    term_id varchar UNIQUE,
    desc_id varchar NOT NULL,
    exp_ck bigint,
    liab_var jsonb,
    asset_vars jsonb
);

CREATE TABLE pool_term_exps (
    exp_ck bigint UNIQUE,
    sup_exp_ck bigint,
    kind smallint,
    spec jsonb
);

CREATE TABLE pool_impl_binds (
    impl_qn ltree UNIQUE,
    impl_id varchar,
    kind smallint
);

CREATE INDEX pool_impl_qn_gist_idx ON pool_impl_binds USING GIST (impl_qn);

CREATE TABLE pool_exec_evals (
    exec_id varchar UNIQUE,
    exec_rn bigint,
    liab_mode smallint
);

CREATE TABLE pool_exec_vars (
    exec_id varchar,
    exec_rn bigint,
    comm_id varchar,
    chnl_id varchar,
    chnl_ph varchar,
    exp_ck bigint,
    side smallint
);

CREATE TABLE pool_struct_vars (
) INHERITS (pool_exec_vars);

CREATE TABLE pool_linear_vars (
) INHERITS (pool_exec_vars);

CREATE TABLE pool_exec_comms (
    comm_id varchar UNIQUE,
    provider_exec_id varchar NOT NULL,
    provider_exec_rn bigint NOT NULL,
    client_exec_id varchar NOT NULL,
    client_exec_rn bigint NOT NULL,
    offset_nr bigint,
    act_nr bigint
);

CREATE TABLE pool_comm_acts (
    comm_id varchar,
    act_nr bigint,
    exec_id varchar,
    exec_rn bigint,
    chnl_id varchar,
    kind smallint,
    exp jsonb
);

CREATE TABLE proc_desc_roots (
    desc_id varchar PRIMARY KEY,
    desc_rn bigint NOT NULL
);

CREATE TABLE proc_desc_binds (
    desc_qn ltree UNIQUE,
    desc_id varchar,
    kind smallint
);

CREATE INDEX proc_desc_qn_gist_idx ON proc_desc_binds USING GIST (desc_qn);

CREATE TABLE proc_type_defs (
    type_id varchar UNIQUE,
    desc_id varchar NOT NULL,
    exp_ck bigint
);

CREATE TABLE proc_type_exps (
    exp_ck bigint UNIQUE,
    sup_exp_ck bigint,
    kind smallint,
    spec jsonb
);

CREATE TABLE proc_desc_terms (
    term_id varchar UNIQUE,
    desc_id varchar NOT NULL,
    exp_ck bigint,
    liab_var jsonb,
    asset_vars jsonb
);

CREATE TABLE proc_term_exps (
    exp_ck bigint UNIQUE,
    sup_exp_ck bigint,
    kind smallint,
    spec jsonb
);

CREATE TABLE proc_impl_binds (
    impl_qn ltree UNIQUE,
    impl_id varchar,
    kind smallint
);

CREATE INDEX proc_impl_qn_gist_idx ON proc_impl_binds USING GIST (impl_qn);

CREATE TABLE proc_exec_evals (
    exec_id varchar UNIQUE,
    exec_rn bigint,
    liab_mode smallint
);

CREATE TABLE proc_exec_vars (
    exec_id varchar,
    exec_rn bigint,
    comm_id varchar,
    chnl_id varchar,
    chnl_ph varchar,
    exp_ck bigint,
    side smallint
);

CREATE TABLE proc_struct_vars (
) INHERITS (proc_exec_vars);

CREATE TABLE proc_linear_vars (
) INHERITS (proc_exec_vars);

CREATE TABLE proc_exec_comms (
    comm_id varchar UNIQUE,
    provider_exec_id varchar NOT NULL,
    provider_exec_rn bigint NOT NULL,
    client_exec_id varchar NOT NULL,
    client_exec_rn bigint NOT NULL,
    offset_nr bigint,
    act_nr bigint
);

CREATE TABLE proc_comm_acts (
    comm_id varchar,
    act_nr bigint,
    exec_id varchar,
    exec_rn bigint,
    chnl_id varchar,
    kind smallint,
    exp jsonb
);
