package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// RBACActions holds the schema definition for the RBACActions entity.
type RBACActions struct {
	ent.Schema
}

// Fields of the RBACActions.
func (RBACActions) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
	}
}

// Edges of the RBACActions.
func (RBACActions) Edges() []ent.Edge {
	return nil
}

func (RBACActions) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "rbac_actions"},
	}
}
