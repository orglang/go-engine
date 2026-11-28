package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// ProgRootMixin defines the fields shared by pool and process Program
// aggregate root tables that are used for aggregate revision tracking.
type ProgRootMixin struct {
	mixin.Schema
}

// Fields of the ProgRootMixin.
func (ProgRootMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("prog_id").Unique().Immutable(),
		field.Int64("prog_rn"),
	}
}
