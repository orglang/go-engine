-- связка описаний с квалифицированными синонимами 
CREATE TABLE pool_desc_binds (
	desc_qn ltree UNIQUE,
	desc_id varchar,
	kind smallint
);

CREATE INDEX pool_desc_qn_gist_idx ON pool_desc_binds USING GIST (desc_qn);

CREATE TABLE pool_type_defs (
	type_id varchar UNIQUE,
	type_rn bigint,
	exp_vk bigint
);

CREATE TABLE pool_type_exps (
	exp_vk bigint UNIQUE,
	sup_exp_vk bigint,
	kind smallint,
	spec jsonb
);

CREATE TABLE pool_term_decs (
	term_id varchar UNIQUE,
	term_rn bigint,
    liab_var jsonb,
    asset_vars jsonb
);

-- связка воплощений с квалифицированными синонимами 
CREATE TABLE pool_impl_binds (
	impl_qn ltree UNIQUE,
	impl_id varchar,
	kind smallint
);

CREATE INDEX pool_impl_qn_gist_idx ON pool_impl_binds USING GIST (impl_qn);

CREATE TABLE pool_comp_execs (
	comp_id varchar UNIQUE,
	comp_rn bigint,
	liab_mode smallint
);

CREATE TABLE pool_comp_vars (
	comp_id varchar,
	comp_rn bigint, -- только для сортировки
	comm_id varchar,
	chnl_id varchar,
	chnl_ph varchar,
	exp_vk bigint,
	side smallint
);

CREATE TABLE pool_struct_vars (
) INHERITS (pool_comp_vars);

CREATE TABLE pool_linear_vars (
) INHERITS (pool_comp_vars);

CREATE TABLE pool_comm_exchs (
	comm_id varchar UNIQUE,
	comm_rn bigint,
	offset_nr bigint
);

CREATE TABLE pool_comm_turns (
	comm_id varchar,
	comm_rn bigint,
	comp_id varchar,
	chnl_id varchar,
	kind smallint,
	exp jsonb
);

-- связка описаний с квалифицированными синонимами 
CREATE TABLE proc_desc_binds (
	desc_qn ltree UNIQUE,
	desc_id varchar,
	kind smallint
);

CREATE INDEX proc_desc_qn_gist_idx ON proc_desc_binds USING GIST (desc_qn);

CREATE TABLE proc_type_defs (
	type_id varchar UNIQUE,
	type_rn bigint,
	exp_vk bigint
);

CREATE TABLE proc_type_exps (
	exp_vk bigint UNIQUE,
	sup_exp_vk bigint,
	kind smallint,
	spec jsonb
);

CREATE TABLE proc_term_decs (
	term_id varchar UNIQUE,
	term_rn bigint,
    liab_var jsonb,
    asset_vars jsonb
);

-- связка воплощений с квалифицированными синонимами 
CREATE TABLE proc_impl_binds (
	impl_qn ltree UNIQUE,
	impl_id varchar,
	kind smallint
);

CREATE INDEX proc_impl_qn_gist_idx ON proc_impl_binds USING GIST (impl_qn);

CREATE TABLE proc_comp_execs (
	comp_id varchar UNIQUE,
	comp_rn bigint,
	liab_mode smallint
);

CREATE TABLE proc_comp_vars (
	comp_id varchar,
	comp_rn bigint, -- только для сортировки
	comm_id varchar,
	chnl_id varchar,
	chnl_ph varchar,
	exp_vk bigint,
	side smallint
);

CREATE TABLE proc_struct_vars (
) INHERITS (proc_comp_vars);

CREATE TABLE proc_linear_vars (
) INHERITS (proc_comp_vars);

CREATE TABLE proc_comm_exchs (
	comm_id varchar UNIQUE,
	comm_rn bigint,
	offset_nr bigint
);

CREATE TABLE proc_comm_turns (
	comm_id varchar,
	comm_rn bigint,
	comp_id varchar,
	chnl_id varchar,
	kind smallint,
	exp jsonb
);
