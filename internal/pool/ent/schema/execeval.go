package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"orglang/go-engine/internal/adt/comproot/adapter/entschema"
)

// ExecEval describes a persisted pool computation execution.
type ExecEval struct {
	ent.Schema
}

// Mixin of the ExecEval.
func (ExecEval) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entschema.CompRootMixin{},
	}
}

// Fields of the ExecEval.
func (ExecEval) Fields() []ent.Field {
	return []ent.Field{
		field.Int16("liab_mode"),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (ExecEval) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pool_exec_evals"},
	}
}
