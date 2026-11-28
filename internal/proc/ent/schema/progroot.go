package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"

	"orglang/go-engine/internal/adt/progroot/adapter/entschema"
)

// ProgRoot describes a persisted process Program aggregate root.
type ProgRoot struct {
	ent.Schema
}

// Mixin of the ProgRoot.
func (ProgRoot) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entschema.ProgRootMixin{},
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (ProgRoot) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proc_desc_roots"},
	}
}
