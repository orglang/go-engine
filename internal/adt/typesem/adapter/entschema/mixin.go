package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// TypeSemMixin defines the identity and owning Description fields shared by
// pool and process Type definition tables.
type TypeSemMixin struct {
	mixin.Schema
}

// Fields of the TypeSemMixin.
func (TypeSemMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("type_id").Immutable(),
		field.String("desc_id"),
	}
}
