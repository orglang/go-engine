package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	termvaradapter "orglang/go-engine/internal/adt/termvar/adapter"
)

type TermDec struct {
	ent.Schema
}

// Fields of the TermDec.
func (TermDec) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("term_id").Immutable(),
		field.String("desc_id"),
		field.Int64("exp_ck").Optional(),
		field.JSON("liab_var", termvaradapter.VarRec{}),
		field.JSON("asset_vars", []termvaradapter.VarRec{}),
	}
}

// Annotations maps the schema to its owning Description's Term identity table.
func (TermDec) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proc_desc_terms"},
	}
}
