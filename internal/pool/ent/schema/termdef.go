package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	termvaradapter "orglang/go-engine/internal/adt/termvar/adapter"
)

// TermDef describes a persisted pool term definition.
type TermDef struct {
	ent.Schema
}

// Fields of the TermDef.
func (TermDef) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("term_id").Immutable(),
		field.String("desc_id"),
		field.JSON("liab_var", termvaradapter.VarRec{}),
		field.JSON("asset_vars", []termvaradapter.VarRec{}),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (TermDef) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pool_term_decs"},
	}
}
