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

// Clinic holds the schema definition for the Clinic entity.
type Clinic struct {
	ent.Schema
}

// Fields of the Clinic.
func (Clinic) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("clinic_id").StructTag(`json:"clinic_id"`),
		field.String("clinic_name"),
		field.Int("user_id").Optional(),
		field.Bool("is_active").StorageKey("isActive").StructTag(`json:"isActive"`),
		field.Int("clinic_account_id").Optional(),
		field.String("clinic_name_old_system").Optional(),
		field.Time("clinic_signup_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now),
		field.Time("clinic_updated_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional().UpdateDefault(time.Now),
	}
}

// Edges of the Clinic.
func (Clinic) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Field("user_id").Ref("clinic").Unique(),
		edge.To("clinic_contacts", Contact.Type),
		edge.To("clinic_addresses", Address.Type),
		//TODO: ref range group
		edge.To("customers", Customer.Type),
		// deprecate
		//edge.To("internal_users", InternalUser.Type),
		edge.To("clinic_settings", Setting.Type),
		edge.To("clinic_orders", OrderInfo.Type),
		edge.To("clinic_patients", Patient.Type),
		//TODO: customer role on clinics? (TBD)
		//TODO: beta program participants
		edge.To("clinic_beta_program_participations", BetaProgramParticipation.Type),
		edge.To("clinic_customer_settings", CustomerSettingOnClinics.Type),
		edge.To("clinic_customer_addresses", CustomerAddressOnClinics.Type),
		edge.To("clinic_customer_contacts", CustomerContactOnClinics.Type),
	}
}

func (Clinic) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "clinic",
		},
	}
}
