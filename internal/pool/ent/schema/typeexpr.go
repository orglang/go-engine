package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// ExprSpec is the JSON payload stored for a pool type-expression state.
type ExprSpec struct {
	Link   string    `json:"link,omitempty"`
	Tensor *ExprProd  `json:"tensor,omitempty"`
	Lolli  *ExprProd  `json:"lolli,omitempty"`
	Plus   *ExprSum   `json:"plus,omitempty"`
	With   *ExprSum   `json:"with,omitempty"`
	Up     *ExprShift `json:"up,omitempty"`
	Down   *ExprShift `json:"down,omitempty"`
}

// ExprProd describes a product expression in a type-expression payload.
type ExprProd struct {
	ValExpCK  int64 `json:"on"`
	ContExpCK int64 `json:"to"`
}

// ExprSum describes a sum expression in a type-expression payload.
type ExprSum struct {
	ProcQNs   []string `json:"on"`
	ContExpCK int64    `json:"to"`
}

// ExprShift describes a shift expression in a type-expression payload.
type ExprShift struct {
	ContExpCK int64 `json:"to"`
}

// TypeExpr describes one persisted state in a pool type expression.
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
		entsql.Annotation{Table: "pool_type_exprs"},
	}
}
