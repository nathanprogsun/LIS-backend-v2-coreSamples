package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Contact holds the schema definition for the Contact entity.
type Contact struct {
	ent.Schema
}

// Fields of the Contact.
func (Contact) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("contact_id").StructTag(`json:"contact_id"`),
		field.String("contact_description").Optional(),
		field.String("contact_details"),
		field.String("contact_type").Optional(),
		field.Bool("is_primary_contact").Optional().Default(false),
		field.Bool("is_2fa_contact").Default(false),
		field.Int("customer_id").Optional(),
		field.Int("patient_id").Optional(),
		field.Int("clinic_id").Optional(),
		field.Int("internal_user_id").Optional(),
		field.Int("user_id").Optional(),
		field.Int("contact_level").Default(1),
		field.String("contact_level_name").Default("Customer"),
		field.Int("group_contact_id").Optional(),
		field.Bool("apply_to_all_group_member").Default(false).
			StorageKey("applyToAllGroupMember").
			StructTag("applyToAllGroupMember"),
		field.Bool("is_group_contact").Default(false).
			StorageKey("isGroupContact").StructTag(`json:"isGroupContact"`),
		field.Bool("use_as_default_create_contact").Default(true).
			StorageKey("useAsDefaultCreateContact").StructTag(`json:"useAsDefaultCreateContact"`),
		field.Bool("use_group_contact").Optional().Default(false).
			StorageKey("useGroupContact").StructTag(`json:"useGroupContact"`),
	}
}

// Edges of the Contact.
func (Contact) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("clinic", Clinic.Type).Ref("clinic_contacts").Field("clinic_id").Unique(),
		edge.From("patient", Patient.Type).Ref("patient_contacts").Field("patient_id").Unique(),
		edge.From("customer", Customer.Type).Ref("customer_contacts").Field("customer_id").Unique(),
		edge.To("customer_clinic_mappings", CustomerContactOnClinics.Type),
		edge.To("group_contact", Contact.Type).Field("group_contact_id").Unique().From("member_contacts"),
		edge.From("internal_user", InternalUser.Type).Ref("internal_user_contacts").Field("internal_user_id").Unique(),
		edge.To("orders", OrderInfo.Type),
	}
}

// Annotations of the Contact.
func (Contact) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "contact",
		},
	}
}
