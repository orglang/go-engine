package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// TypeSemMixin defines the fields shared by pool and process type definition
// tables that are used for ADT semantic revision tracking.
type TypeSemMixin struct {
	mixin.Schema
}

// Fields of the TypeSemMixin.
func (TypeSemMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("type_id").Immutable(),
		field.Int64("type_rn"),
	}
}
