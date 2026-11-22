package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	termvaradapter "orglang/go-engine/internal/adt/termvar/adapter"
)

// TermDec describes a persisted process term declaration.
type TermDec struct {
	ent.Schema
}

// Fields of the TermDec.
func (TermDec) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("term_id").Immutable(),
		field.Int64("term_rn"),
		field.JSON("liab_var", termvaradapter.VarRec{}),
		field.JSON("asset_vars", []termvaradapter.VarRec{}),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (TermDec) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proc_term_decs"},
	}
}
