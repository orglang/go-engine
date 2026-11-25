package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	termvaradapter "orglang/go-engine/internal/adt/termvar/adapter"
)

// Deprecated: TermDef is the legacy declaration-shaped model for a Term. Use the Term identity and Definition repositories.
type TermDef struct {
	ent.Schema
}

// Fields of the TermDef.
func (TermDef) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("term_id").Immutable(),
		field.String("desc_id"),
		field.Int64("exp_ck").Optional(),
		field.JSON("liab_var", termvaradapter.VarRec{}),
		field.JSON("asset_vars", []termvaradapter.VarRec{}),
	}
}

// Annotations maps the schema to its owning Description's Term identity table.
func (TermDef) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pool_desc_terms"},
	}
}
