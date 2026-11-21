package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// TypeExp describes one persisted state in a pool type expression.
type TypeExp struct {
	ent.Schema
}

// Fields of the TypeExp.
func (TypeExp) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("exp_ck").Unique().Immutable(),
		field.Int64("sup_exp_ck"),
		field.Int16("kind"),
		field.JSON("spec", json.RawMessage{}).Optional(),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (TypeExp) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pool_type_exps"},
	}
}
