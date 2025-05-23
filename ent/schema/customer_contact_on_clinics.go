package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CustomerContactOnClinics holds the schema definition for the relation table.
type CustomerContactOnClinics struct {
	ent.Schema
}

// Fields of the CustomerContactOnClinics.
func (CustomerContactOnClinics) Fields() []ent.Field {
	return []ent.Field{
		field.Int("customer_id"),
		field.Int("clinic_id"),
		field.Int("contact_id"),
		field.String("contact_type").Optional(),
	}
}

// Edges of the CustomerContactOnClinics.
func (CustomerContactOnClinics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("customer", Customer.Type).
			Ref("customer_contacts_on_clinics").
			Field("customer_id").
			Unique().
			Required(),

		edge.From("clinic", Clinic.Type).
			Ref("clinic_customer_contacts").
			Field("clinic_id").
			Unique().
			Required(),

		edge.From("contact", Contact.Type).
			Ref("customer_clinic_mappings").
			Field("contact_id").
			Unique().
			Required(),
	}
}

// Indexes defines unique constraints.
func (CustomerContactOnClinics) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("customer_id", "clinic_id", "contact_id").Unique(),
	}
}

// Annotations of the CustomerContactOnClinics.
func (CustomerContactOnClinics) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "customer_contact_on_clinics",
		},
	}
}
