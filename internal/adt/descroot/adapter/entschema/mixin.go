package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// DescRootMixin defines the fields shared by pool and process Description
// aggregate root tables that are used for aggregate revision tracking.
type DescRootMixin struct {
	mixin.Schema
}

// Fields of the DescRootMixin.
func (DescRootMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("desc_id").Unique().Immutable(),
		field.Int64("desc_rn"),
	}
}
