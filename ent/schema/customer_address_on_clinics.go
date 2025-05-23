package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CustomerAddressOnClinics holds the schema definition for the relation table.
type CustomerAddressOnClinics struct {
	ent.Schema
}

// Fields of the CustomerAddressOnClinics.
func (CustomerAddressOnClinics) Fields() []ent.Field {
	return []ent.Field{
		field.Int("customer_id"),
		field.Int("clinic_id"),
		field.Int("address_id"),
		field.String("address_type").Optional(),
	}
}

// Edges of the CustomerAddressOnClinics.
func (CustomerAddressOnClinics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("customer", Customer.Type).
			Ref("customer_addresses_on_clinics").
			Field("customer_id").
			Unique().
			Required(),

		edge.From("clinic", Clinic.Type).
			Ref("clinic_customer_addresses").
			Field("clinic_id").
			Unique().
			Required(),

		edge.From("address", Address.Type).
			Ref("customer_clinic_mappings").
			Field("address_id").
			Unique().
			Required(),
	}
}

// Indexes defines unique constraints.
func (CustomerAddressOnClinics) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("customer_id", "clinic_id", "address_id").Unique(),
	}
}

// Annotations of the CustomerAddressOnClinics.
func (CustomerAddressOnClinics) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "customer_address_on_clinics",
		},
	}
}
