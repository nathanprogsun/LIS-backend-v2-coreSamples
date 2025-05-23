package schema

import (
	"time"

	"entgo.io/ent/dialect"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Customer holds the schema definition for the Customer entity.
type Customer struct {
	ent.Schema
}

// Fields of the Customer.
func (Customer) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("customer_id").StructTag(`json:"customer_id"`),
		field.Int("user_id").Optional().Unique(),
		field.String("customer_type").Default("Normal Customer"),
		field.String("customer_first_name").Optional(),
		field.String("customer_last_name").Optional(),
		field.String("customer_middle_name").Optional(),
		field.String("customer_type_id").Optional(),
		field.String("customer_suffix").Optional(),
		field.String("customer_samples_received").Optional(),
		field.Time("customer_request_submit_time").Optional(),
		field.Time("customer_signup_time").Optional().Default(time.Now),
		field.Bool("is_active").
			StorageKey("isActive").StructTag(`json:"isActive"`).Default(true),
		field.Int("sales_id").Optional(),
		field.String("customer_npi_number").Optional(),
		field.String("referral_source").Optional(),
		field.Bool("order_placement_allowed").Optional().Default(true),
		field.Bool("beta_program_enabled").Optional().Default(false),
		field.Time("onboarding_questionnaire_filled_on").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional(),
	}
}

// Edges of the Customer.
func (Customer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("samples", Sample.Type),
		edge.To("customer_contacts", Contact.Type),
		edge.To("customer_addresses", Address.Type),
		edge.From("clinics", Clinic.Type).Ref("customers"),
		//former internal_user @relation(ActiveSales)
		edge.From("sales", InternalUser.Type).Ref("customers").Field("sales_id").Unique(),
		edge.From("user", User.Type).Ref("customer").Field("user_id").Unique(),
		//TODO: customer_internal_notes
		edge.To("orders", OrderInfo.Type),
		edge.To("current_patients", Patient.Type),
		edge.To("patients", Patient.Type),
		//TODO: customer_reference_range
		//TODO: customer_interests
		edge.To("customer_beta_program_participations", BetaProgramParticipation.Type),
		edge.To("customer_settings_on_clinics", CustomerSettingOnClinics.Type),
		edge.To("customer_addresses_on_clinics", CustomerAddressOnClinics.Type),
		edge.To("customer_contacts_on_clinics", CustomerContactOnClinics.Type),
	}
}

// Annotations of the Customer.
func (Customer) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "customer",
		},
	}
}
