CREATE TABLE type_defs (
	def_id varchar(36),
	def_rn bigint,
	title varchar(64)
);

CREATE TABLE type_exps (
	exp_id varchar(36),
	from_id varchar(36),
	kind smallint,
	spec jsonb
);

CREATE TABLE proc_decs (
	dec_id varchar(36),
	dec_rn bigint,
	title text
);

CREATE TABLE dec_pes (
	dec_id varchar(36),
	chnl_ph varchar(64),
	type_qn ltree,
	from_rn bigint,
	to_rn bigint
);

CREATE TABLE dec_ces (
	dec_id varchar(36),
	chnl_ph varchar(64),
	type_qn ltree,
	from_rn bigint,
	to_rn bigint
);

CREATE TABLE dec_subs (
	dec_id varchar(36),
	dec_qn ltree,
	from_rn bigint,
	to_rn bigint
);

CREATE TABLE pool_execs (
	exec_id varchar(36),
	exec_rn bigint,
	title varchar(64),
	proc_id varchar(36),
	sup_exec_id varchar(36)
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

CREATE TABLE syn_decs (
	dec_id varchar(36),
	dec_qn ltree UNIQUE,
	from_rn bigint,
	to_rn bigint,
	kind smallint
);

CREATE INDEX sym_gist_idx ON syn_decs USING GIST (sym);
