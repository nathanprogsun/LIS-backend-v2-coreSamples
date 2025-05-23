package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("user_id").StructTag(`json:"user_id"`),
		field.String("user_name").StorageKey("username").StructTag("username").Unique(),
		field.String("email_user_id").Unique().Optional(),
		field.String("password"),
		field.String("two_factor_authentication_secret").
			StorageKey("twoFactorAuthenticationSecret").Optional(),
		field.Bool("is_two_factor_authentication_enabled").
			StorageKey("isTwoFactorAuthenticationEnabled").Default(false),
		field.String("user_group").Optional(),
		field.Bool("imported_user_with_salt_password").Default(false),
		field.Bool("is_active").StorageKey("isActive").StructTag(`json:"isActive"`).Default(true),
		// user_permission is deprecated for not being used
		// isInternalAdminUser, hasAdminPanelAccess, hasLISAccess are deprecated for now unless
		// we have trouble later migrating to the new RBAC system
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("clinic", Clinic.Type),
		edge.To("customer", Customer.Type),
		edge.To("patient", Patient.Type),
		edge.To("internal_user", InternalUser.Type),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user"},
	}
}
