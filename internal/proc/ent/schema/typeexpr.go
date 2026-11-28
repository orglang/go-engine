package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// ExprSpec is the JSON payload stored for a process type-expression state.
type ExprSpec struct {
	Link   string   `json:"link,omitempty"`
	Tensor *ExprProd `json:"tensor,omitempty"`
	Lolli  *ExprProd `json:"lolli,omitempty"`
	Plus   []ExprSum `json:"plus,omitempty"`
	With   []ExprSum `json:"with,omitempty"`
}

// ExprProd describes a product expression in a type-expression payload.
type ExprProd struct {
	ValExpCK  int64 `json:"on"`
	ContExpCK int64 `json:"to"`
}

// ExprSum describes a choice branch in a type-expression payload.
type ExprSum struct {
	LabQN     string `json:"on"`
	ContExpCK int64  `json:"to"`
}

// TypeExpr describes one persisted state in a process type expression.
type TypeExpr struct {
	ent.Schema
}

// Fields of the TypeExpr.
func (TypeExpr) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("expr_ck").Unique().Immutable(),
		field.Int64("sup_expr_ck"),
		field.Int16("kind"),
		field.JSON("spec", ExprSpec{}).Optional(),
	}
}

// Annotations maps the schema to its existing PostgreSQL table.
func (TypeExpr) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proc_type_exprs"},
	}
}
