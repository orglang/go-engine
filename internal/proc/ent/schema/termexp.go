package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// TermExp describes one persisted process term-definition expression node.
type TermExp struct {
	ent.Schema
}

// Fields of the TermExp.
func (TermExp) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("exp_ck").Unique().Immutable(),
		field.Int64("sup_exp_ck"),
		field.Int16("kind"),
		field.JSON("spec", json.RawMessage{}).Optional(),
	}
}

// Annotations maps the schema to its PostgreSQL expression table.
func (TermExp) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proc_term_exps"},
	}
}
