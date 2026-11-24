package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	termexpadapter "orglang/go-engine/internal/proc/termexp/adapter"
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
		field.JSON("spec", termexpadapter.StateSpec{}).Optional(),
	}
}

// Annotations maps the schema to its PostgreSQL expression table.
func (TermExp) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proc_term_exps"},
	}
}
