package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"

	"orglang/go-engine/internal/adt/descroot/adapter/entschema"
)

// DescRoot describes a persisted pool Description aggregate root.
type DescRoot struct {
	ent.Schema
}

// Mixin of the DescRoot.
func (DescRoot) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entschema.DescRootMixin{},
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (DescRoot) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pool_desc_roots"},
	}
}
