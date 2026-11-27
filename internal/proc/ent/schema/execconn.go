package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// ExecConn describes a persisted process computation connection.
type ExecConn struct {
	ent.Schema
}

// Fields of the ExecConn.
func (ExecConn) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("conn_id").Unique().Immutable(),
		field.String("provider_exec_id"),
		field.Int64("provider_exec_rn"),
		field.String("client_exec_id"),
		field.Int64("client_exec_rn"),
		field.Int64("offset_nr").Optional(),
		field.Int64("act_nr").Optional(),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (ExecConn) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proc_exec_conns"},
	}
}
