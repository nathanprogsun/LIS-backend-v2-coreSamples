package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// RBACResources holds the schema definition for the RBACResources entity.
type RBACResources struct {
	ent.Schema
}

// Fields of the RBACResources.
func (RBACResources) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
		field.String("description").Default("N/A"),
	}
}

// Edges of the RBACResources.
func (RBACResources) Edges() []ent.Edge {
	return nil
}

func (RBACResources) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "rbac_resources"},
	}
}
