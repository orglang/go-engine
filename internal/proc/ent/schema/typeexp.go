package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// ExpSpec is the JSON payload stored for a process type-expression state.
type ExpSpec struct {
	Link   string    `json:"link,omitempty"`
	Tensor *ExpProd  `json:"tensor,omitempty"`
	Lolli  *ExpProd  `json:"lolli,omitempty"`
	Plus   []ExpSum  `json:"plus,omitempty"`
	With   []ExpSum  `json:"with,omitempty"`
}

// ExpProd describes a product expression in a type-expression payload.
type ExpProd struct {
	ValExpCK  int64 `json:"on"`
	ContExpCK int64 `json:"to"`
}

// ExpSum describes a choice branch in a type-expression payload.
type ExpSum struct {
	LabQN     string `json:"on"`
	ContExpCK int64  `json:"to"`
}

// TypeExp describes one persisted state in a process type expression.
type TypeExp struct {
	ent.Schema
}

// Fields of the TypeExp.
func (TypeExp) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("exp_ck").Unique().Immutable(),
		field.Int64("sup_exp_ck"),
		field.Int16("kind"),
		field.JSON("spec", ExpSpec{}).Optional(),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (TypeExp) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proc_type_exps"},
	}
}
