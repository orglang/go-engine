package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// ExecRootMixin defines the fields shared by pool and process execution
// tables that are used for ADT semantic revision tracking.
type ExecRootMixin struct {
	mixin.Schema
}

// Fields of the ExecRootMixin.
func (ExecRootMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("exec_id").Unique().Immutable(),
		field.Int64("root_rn").StorageKey("exec_rn"),
	}
}
