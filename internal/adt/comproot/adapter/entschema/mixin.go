package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// CompRootMixin defines the fields shared by pool and process computation
// execution tables that are used for ADT semantic revision tracking.
type CompRootMixin struct {
	mixin.Schema
}

// Fields of the CompRootMixin.
func (CompRootMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("comp_id").Unique().Immutable(),
		field.Int64("comp_rn"),
	}
}
