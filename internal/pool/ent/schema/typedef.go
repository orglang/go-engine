package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"

	"orglang/go-engine/internal/adt/typesem/adapter/entschema"
)

// TypeDef describes a persisted pool type definition.
type TypeDef struct {
	ent.Schema
}

// Mixin of the TypeDef.
func (TypeDef) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entschema.TypeSemMixin{},
	}
}

// Fields of the TypeDef.
func (TypeDef) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("exp_ck"),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (TypeDef) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pool_type_defs"},
	}
}
