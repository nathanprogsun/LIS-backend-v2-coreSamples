package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CustomerSettingOnClinics holds the schema definition for the CustomerSettingOnClinics entity.
type CustomerSettingOnClinics struct {
	ent.Schema
}

// Fields of the CustomerSettingOnClinics.
func (CustomerSettingOnClinics) Fields() []ent.Field {
	return []ent.Field{
		field.Int("customer_id"),
		field.Int("clinic_id"),
		field.Int("setting_id"),
		field.String("setting_name"),
	}
}

// Edges of the CustomerSettingOnClinics.
func (CustomerSettingOnClinics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("customer", Customer.Type).
			Ref("customer_settings_on_clinics").
			Unique().
			Field("customer_id").
			Required(),
		edge.From("clinic", Clinic.Type).
			Ref("clinic_customer_settings").
			Unique().
			Field("clinic_id").
			Required(),
		edge.From("setting", Setting.Type).
			Ref("clinic_customers").
			Unique().
			Field("setting_id").
			Required(),
	}
}

func (CustomerSettingOnClinics) Indexes() []ent.Index {
	return []ent.Index{
		// Unique constraint for (customer_id, clinic_id, setting_id)
		index.Fields("customer_id", "clinic_id", "setting_id").
			Unique(),

		// Unique constraint for (customer_id, clinic_id, setting_name)
		index.Fields("customer_id", "clinic_id", "setting_name").
			Unique(),

		// Custom named index for (customer_id, clinic_id, setting_id)
		index.Fields("customer_id", "clinic_id", "setting_id").
			StorageKey("customer_setting_on_clinics_customer_id_clinic_id_setting_id_key"),
	}
}

// Annotations of the CustomerSettingOnClinics.
func (CustomerSettingOnClinics) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "customer_setting_on_clinics",
		},
	}
}
