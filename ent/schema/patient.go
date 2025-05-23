package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Patient holds the schema definition for the Patient entity.
type Patient struct {
	ent.Schema
}

// Fields of the Patient.
func (Patient) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("patient_id").StructTag(`json:"patient_id"`),
		field.Int("user_id").Optional(),
		field.String("patient_type").Default("Standalone Patient"),
		field.Int("original_patient_id").Optional(),
		field.String("patient_gender").Optional(),
		field.String("patient_first_name").Optional(),
		field.String("patient_last_name").Optional(),
		field.String("patient_middle_name").Optional(),
		field.String("patient_medical_record_number").Optional(),
		field.String("patient_legal_firstname").Optional(),
		field.String("patient_legal_lastname").Optional(),
		field.String("patient_honorific").Optional(),
		field.String("patient_suffix").Optional(),
		field.String("patient_marital").Optional(),
		field.String("patient_ethnicity").Optional(),
		field.String("patient_birthdate").Optional(),
		field.String("patient_ssn").Optional(),
		field.String("patient_height").Optional(),
		field.String("patient_weight").Optional(),
		field.Int("officeally_id").Optional(),
		field.String("patient_ny_waive_form_issue_status").Default("no_ny_waive_form_issue"),
		field.Time("patient_create_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional().Default(time.Now),
		field.Int("customer_id").Optional(),
		field.Bool("isActive").StorageKey("isActive").Default(true),
		field.Bool("patient_flagged").Default(false),
		field.Time("patient_service_date").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional(),
		field.String("patient_description").Optional(),
		field.String("patient_language").Optional().Default("English"),
	}
}

// Edges of the Patient.
func (Patient) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("current_customer", Customer.Type).Field("customer_id").Unique().Ref("current_patients"),
		edge.From("patient_customers", Customer.Type).Ref("patients"),
		edge.To("samples", Sample.Type),
		edge.To("patient_contacts", Contact.Type),
		edge.To("patient_addresses", Address.Type),
		edge.From("patient_clinics", Clinic.Type).Ref("clinic_patients"),
		edge.From("user", User.Type).Field("user_id").Unique().Ref("patient"),
		edge.To("patient_weight_height_history", PatientWeightHeight.Type),
		//TODO:internal notes
		//edge.To("patient_orders", OrderInfo.Type),
		edge.To("patient_settings", Setting.Type),
	}
}

func (Patient) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "patient"},
	}
}
