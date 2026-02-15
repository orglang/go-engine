CREATE TABLE type_defs (
	def_id varchar(36),
	def_rn bigint,
	syn_vk bigint,
	exp_vk bigint
);

CREATE TABLE type_exps (
	exp_vk bigint UNIQUE,
	sup_exp_vk bigint,
	def_id varchar(36),
	def_rn bigint,
	kind smallint,
	spec jsonb
);

CREATE TABLE xact_defs (
	def_id varchar(36),
	def_rn bigint,
	syn_vk bigint,
	exp_vk bigint
);

CREATE TABLE xact_exps (
	exp_vk bigint UNIQUE,
	sup_exp_vk bigint,
	def_id varchar(36),
	def_rn bigint,
	kind smallint,
	spec jsonb
);

CREATE TABLE pool_decs (
	dec_id varchar(36),
	dec_rn bigint,
	syn_vk bigint,
    client_brs jsonb,
    provider_br jsonb
);

CREATE TABLE pool_execs (
	exec_id varchar(36),
	exec_rn bigint,
	proc_id varchar(36)
);

CREATE TABLE pool_caps (
	pool_id varchar(36),
	sig_id varchar(36),
	rev bigint
);

CREATE TABLE pool_deps (
	pool_id varchar(36),
	sig_id varchar(36),
	rev bigint
);

-- передачи каналов (провайдерская сторона)
-- по истории передач определяем текущего провайдера
CREATE TABLE pool_liabs (
	proc_id varchar(36),
	pool_id varchar(36),
	rev bigint
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
	chnl_ph varchar(36),
	chnl_id varchar(36),
	state_id varchar(36),
	exec_rn bigint
);

CREATE TABLE proc_steps (
	exec_id varchar(36),
	exec_rn bigint,
	chnl_id varchar(36),
	kind smallint,
	proc_er jsonb
);

CREATE TABLE pool_sups (
	pool_id varchar(36),
	sup_pool_id varchar(36),
	rev bigint
);

CREATE TABLE synonyms (
	syn_qn ltree UNIQUE,
	syn_vk bigint UNIQUE,
	kind smallint
);

CREATE INDEX syn_qn_gist_idx ON synonyms USING GIST (syn_qn);
