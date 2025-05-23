package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Address holds the schema definition for the Address entity.
type Address struct {
	ent.Schema
}

// Fields of the Address.
func (Address) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("address_id").StructTag(`json:"address_id"`),
		field.String("address_type").Optional(),
		field.String("street_address").Optional(),
		field.String("apt_po").Optional(),
		field.String("city").Optional(),
		field.String("state").Optional(),
		field.String("zipcode").Optional(),
		field.String("country").Optional(),
		field.Bool("address_confirmed").Default(true),
		field.Bool("is_primary_address").Default(false),
		field.Int("customer_id").Optional(),
		field.Int("patient_id").Optional(),
		field.Int("clinic_id").Optional(),
		field.Int("internal_user_id").Optional(),
		field.Int("address_level").Default(1),
		field.String("address_level_name").Default("Customer"),
		field.Bool("apply_to_all_group_member").Default(false).
			StorageKey("applyToAllGroupMember").StructTag(`json:"applyToAllGroupMember"`),
		field.Int("group_address_id").Optional(),
		field.Bool("is_group_address").Default(false).
			StorageKey("isGroupAddress").StructTag(`json:"isGroupAddress"`),
		field.Bool("use_as_default_create_address").Default(true).
			StorageKey("useAsDefaultCreateAddress").StructTag(`json:"useAsDefaultCreateAddress"`),
		field.Bool("use_group_address").Default(false).
			StorageKey("useGroupAddress").StructTag(`json:"useGroupAddress"`),
	}
}

// Edges of the Address.
func (Address) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("clinic", Clinic.Type).Ref("clinic_addresses").Field("clinic_id").Unique(),
		edge.From("customer", Customer.Type).Ref("customer_addresses").Field("customer_id").Unique(),
		edge.To("customer_clinic_mappings", CustomerAddressOnClinics.Type),
		edge.To("group_address", Address.Type).Field("group_address_id").Unique().From("member_addresses"),
		edge.From("internal_user", InternalUser.Type).Ref("internal_user_addresses").Field("internal_user_id").Unique(),
		edge.From("patient", Patient.Type).Ref("patient_addresses").Field("patient_id").Unique(),
		edge.To("orders", OrderInfo.Type),
	}
}

func (Address) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "address",
		},
	}
}
