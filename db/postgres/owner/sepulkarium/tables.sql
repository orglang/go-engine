CREATE TABLE desc_sems (
	desc_id varchar(36) UNIQUE,
	desc_rn bigint,
	kind smallint
);

-- связка описаний с квалифицированными синонимами 
CREATE TABLE desc_binds (
	desc_qn ltree UNIQUE,
	desc_id varchar(36)
);

CREATE INDEX desc_qn_gist_idx ON desc_binds USING GIST (desc_qn);

CREATE TABLE impl_sems (
	impl_id varchar(36) UNIQUE,
	impl_rn bigint,
	kind smallint
);

-- связка воплощений с квалифицированными синонимами 
CREATE TABLE impl_binds (
	impl_qn ltree UNIQUE,
	impl_id varchar(36)
);

CREATE INDEX impl_qn_gist_idx ON impl_binds USING GIST (impl_qn);

CREATE TABLE xact_defs (
	desc_id varchar(36) UNIQUE,
	exp_vk bigint
);

CREATE TABLE xact_exps (
	exp_vk bigint UNIQUE,
	sup_exp_vk bigint,
	desc_id varchar(36),
	desc_rn bigint,
	kind smallint,
	spec jsonb
);

CREATE TABLE pool_decs (
	desc_id varchar(36) UNIQUE,
    provider_vr jsonb,
    client_vrs jsonb
);

CREATE TABLE pool_execs (
	impl_id varchar(36) UNIQUE,
	impl_rn bigint,
	chnl_id varchar(36),
	chnl_ph varchar,
	exp_vk bigint
);

CREATE TABLE pool_vars (
	impl_id varchar(36),
	impl_rn bigint,
	chnl_id varchar(36),
	chnl_ph varchar,
	exp_vk bigint
);

CREATE TABLE type_defs (
	desc_id varchar(36) UNIQUE,
	exp_vk bigint
);

CREATE TABLE type_exps (
	exp_vk bigint UNIQUE,
	sup_exp_vk bigint,
	desc_id varchar(36),
	desc_rn bigint,
	kind smallint,
	spec jsonb
);

CREATE TABLE proc_decs (
	desc_id varchar(36) UNIQUE,
    provider_vr jsonb,
    client_vrs jsonb
);

CREATE TABLE proc_execs (
	impl_id varchar(36) UNIQUE,
	provider_ph varchar
);

CREATE TABLE proc_vars (
	chnl_id varchar(36),
	chnl_ph varchar,
	exp_vk bigint,
	impl_id varchar(36),
	impl_rn bigint
);

CREATE TABLE proc_steps (
	impl_id varchar(36),
	impl_rn bigint,
	chnl_id varchar(36),
	kind smallint,
	proc_er jsonb
);
