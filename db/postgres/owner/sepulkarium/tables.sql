CREATE TABLE desc_execs (
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
    client_brs jsonb,
    provider_br jsonb
);

CREATE TABLE pool_execs (
	exec_id varchar(36),
	exec_rn bigint,
	proc_id varchar(36)
);

CREATE TABLE pool_caps (
	desc_id varchar(36),
	sig_id varchar(36),
	rev bigint
);

CREATE TABLE pool_deps (
	desc_id varchar(36),
	sig_id varchar(36),
	rev bigint
);

-- передачи каналов (провайдерская сторона)
-- по истории передач определяем текущего провайдера
CREATE TABLE pool_liabs (
	proc_id varchar(36),
	desc_id varchar(36),
	rev bigint
);

CREATE TABLE type_defs (
	desc_id varchar(36),
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
	dec_id varchar(36),
	dec_rn bigint,
	syn_vk bigint
);

CREATE TABLE proc_execs (
	exec_id varchar(36),
	exec_rn bigint
);

-- подстановки каналов в процесс
CREATE TABLE proc_binds (
	exec_id varchar(36),
	exec_rn bigint,
	chnl_ph varchar(36),
	chnl_id varchar(36),
	state_id varchar(36)
);

CREATE TABLE proc_steps (
	exec_id varchar(36),
	exec_rn bigint,
	chnl_id varchar(36),
	kind smallint,
	proc_er jsonb
);

CREATE TABLE pool_sups (
	desc_id varchar(36),
	sup_pool_id varchar(36),
	rev bigint
);
