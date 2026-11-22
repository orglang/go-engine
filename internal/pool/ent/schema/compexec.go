package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"

	"orglang/go-engine/internal/adt/compsem/adapter/entschema"
)

// CompExec describes a persisted pool computation execution.
type CompExec struct {
	ent.Schema
}

// Mixin of the CompExec.
func (CompExec) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entschema.CompSemMixin{},
	}
}

// Fields of the CompExec.
func (CompExec) Fields() []ent.Field {
	return []ent.Field{
		field.Int16("liab_mode"),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (CompExec) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pool_comp_execs"},
	}
}
