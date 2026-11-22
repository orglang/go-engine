package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// CompConn describes a persisted process computation connection.
type CompConn struct {
	ent.Schema
}

// Fields of the CompConn.
func (CompConn) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("conn_id").Unique().Immutable(),
		field.String("provider_comp_id"),
		field.Int64("provider_comp_rn"),
		field.String("client_comp_id"),
		field.Int64("client_comp_rn"),
		field.Int64("offset_nr").Optional(),
		field.Int64("act_nr").Optional(),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (CompConn) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proc_comp_conns"},
	}
}
