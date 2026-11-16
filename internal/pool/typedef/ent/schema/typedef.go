package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// TypeDef describes a persisted pool type definition.
type TypeDef struct {
	ent.Schema
}

// Fields of the TypeDef.
func (TypeDef) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("type_id").Immutable(),
		field.Int64("type_rn"),
		field.Int64("exp_ck"),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (TypeDef) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pool_type_defs"},
	}
}
