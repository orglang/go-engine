package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type TypeDef struct {
	ent.Schema
}

// Mixin of the TypeDef.
// Fields of the TypeDef.
func (TypeDef) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("type_id").Immutable(),
		field.String("desc_id"),
		field.Int64("exp_ck"),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (TypeDef) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pool_type_defs"},
	}
}
