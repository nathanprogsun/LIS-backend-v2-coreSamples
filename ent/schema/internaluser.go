package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// InternalUser holds the schema definition for the InternalUser entity.
type InternalUser struct {
	ent.Schema
}

// Fields of the InternalUser.
func (InternalUser) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("internal_user_id").StructTag(`json:"internal_user_id"`),
		field.String("internal_user_role"),
		field.String("internal_user_name").Optional(),
		field.String("internal_user_firstname").Optional(),
		field.String("internal_user_lastname").Optional(),
		field.String("internal_user_middle_name").Optional(),
		field.Bool("internal_user_is_full_time").Default(true),
		field.String("internal_user_email").Optional(),
		field.String("internal_user_phone").Optional(),
		field.Bool("isActive").StorageKey("isActive").StructTag(`json:"isActive"`).Default(true),
		field.Int("user_id").Optional().Unique(),
		field.String("internal_user_type").Optional(),
		field.Int("internal_user_role_id").Optional(),
		// username is deprecated for being duplicated in user
		// internal_user_type_id, internal_user_region is for old system so will be deprecated for now
		// internal_user_permission should be replaced by RBAC system
		// internal_user_department has been deprecated in v1
		// we probably could also remove email and phone for it's saved in contact as well,
		// but it seems to be here for convenience according to Zichen
	}
}

// Edges of the InternalUser.
func (InternalUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sales_team", SalesTeam.Type).Ref("internal_user"),
		edge.To("internal_user_contacts", Contact.Type),
		edge.To("internal_user_addresses", Address.Type),
		//deprecate
		//edge.From("clinics", Clinic.Type).Ref("internal_users"),
		edge.To("customers", Customer.Type),
		edge.From("user", User.Type).Ref("internal_user").Field("user_id").Unique(),
		edge.To("internal_user_settings", Setting.Type),
	}
}

func (InternalUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "internal_user"},
	}
}
