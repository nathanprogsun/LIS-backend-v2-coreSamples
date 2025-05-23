package schema

import (
	"coresamples/model"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// RBACRoles holds the schema definition for the RBACRoles entity.
type RBACRoles struct {
	ent.Schema
}

// Fields of the RBACRoles.
func (RBACRoles) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("internal_name").Unique().NotEmpty(),
		field.Enum("type").Values(model.GetRoleTypes()...),
		field.Int32("clinic_id"),
	}
}

// Edges of the RBACRoles.
func (RBACRoles) Edges() []ent.Edge {
	return nil
}

func (RBACRoles) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "rbac_roles"},
	}
}
