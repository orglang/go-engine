package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	compsemSchema "orglang/go-engine/internal/adt/compsem/adapter/entschema"
)

// CompExec describes a persisted process computation execution.
type CompExec struct {
	ent.Schema
}

// Mixin of the CompExec.
func (CompExec) Mixin() []ent.Mixin {
	return []ent.Mixin{
		compsemSchema.CompSemMixin{},
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
		entsql.Annotation{Table: "proc_comp_execs"},
	}
}
