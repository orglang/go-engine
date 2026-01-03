CREATE TABLE type_def_roots (
	def_id varchar(36),
	title varchar(64),
	def_rn bigint
);

CREATE TABLE type_def_states (
	term_id varchar(36),
	from_id varchar(36),
	kind smallint,
	spec jsonb
);

CREATE TABLE role_subs (
	role_id varchar(36),
	role_fqn ltree,
	from_rn bigint,
	to_rn bigint
);

CREATE TABLE sig_roots (
	sig_id varchar(36),
	title text,
	rev bigint
);

CREATE TABLE sig_pes (
	sig_id varchar(36),
	chnl_key varchar(64),
	role_fqn ltree,
	from_rn bigint,
	to_rn bigint
);

CREATE TABLE sig_ces (
	sig_id varchar(36),
	chnl_key varchar(64),
	role_fqn ltree,
	from_rn bigint,
	to_rn bigint
);

CREATE TABLE sig_subs (
	sig_id varchar(36),
	sig_fqn ltree,
	from_rn bigint,
	to_rn bigint
);

CREATE TABLE pool_roots (
	pool_id varchar(36),
	title varchar(64),
	proc_id varchar(36),
	sup_pool_id varchar(36),
	rev integer
);

CREATE TABLE pool_caps (
	pool_id varchar(36),
	sig_id varchar(36),
	rev integer
);

CREATE TABLE pool_deps (
	pool_id varchar(36),
	sig_id varchar(36),
	rev integer
);

-- передачи каналов (провайдерская сторона)
-- по истории передач определяем текущего провайдера
CREATE TABLE pool_liabs (
	proc_id varchar(36),
	pool_id varchar(36),
	rev integer
);

-- подстановки каналов в процесс
CREATE TABLE proc_bnds (
	proc_id varchar(36),
	chnl_ph varchar(36),
	chnl_id varchar(36),
	state_id varchar(36),
	rev integer
);

CREATE TABLE proc_steps (
	proc_id varchar(36),
	chnl_id varchar(36),
	kind smallint,
	spec jsonb,
	rev integer
);

CREATE TABLE pool_sups (
	pool_id varchar(36),
	sup_pool_id varchar(36),
	rev integer
);

CREATE TABLE steps (
	id varchar(36),
	kind smallint,
	pid varchar(36),
	vid varchar(36),
	spec jsonb
);

CREATE TABLE aliases (
	dec_id varchar(36),
	dec_qn ltree UNIQUE,
	from_rn bigint,
	to_rn bigint,
	kind smallint
);

CREATE INDEX sym_gist_idx ON aliases USING GIST (sym);
