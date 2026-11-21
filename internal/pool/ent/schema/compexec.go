package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// CompExec describes a persisted pool computation execution.
type CompExec struct {
	ent.Schema
}

// Fields of the CompExec.
func (CompExec) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("comp_id").Unique().Immutable(),
		field.Int64("comp_rn"),
		field.Int16("liab_mode"),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (CompExec) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pool_comp_execs"},
	}
}
