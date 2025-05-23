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

// Setting holds the schema definition for the Setting entity.
type Setting struct {
	ent.Schema
}

// Fields of the Setting.
func (Setting) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StorageKey("setting_id").
			StructTag(`json:"setting_id"`).
			Unique().
			Immutable(),
		field.String("setting_name"),
		field.String("setting_group").
			Default("default_customer_setting_group"),
		field.String("setting_description"),
		field.String("setting_value").
			Optional().
			SchemaType(map[string]string{dialect.MySQL: "MEDIUMTEXT"}),
		field.String("setting_type"),
		field.Time("setting_value_updated_time").
			Default(time.Now),
		field.Bool("is_active").
			StorageKey("isActive").
			Default(true).
			StructTag(`json:"isActive"`),
		field.Bool("apply_to_all_group_member").
			Default(false).
			StructTag(`json:"applyToAllGroupMember"`),
		field.Bool("is_official").
			Default(false).
			StructTag(`json:"isOfficial"`),
		field.Int("setting_level").
			Default(1),
		field.String("setting_level_name").
			Default("Customer"),
		field.Bool("use_group_setting").
			Default(false).
			StructTag(`json:"useGroupSetting"`),
	}
}

// Edges of the Setting.
func (Setting) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("clinics", Clinic.Type).Ref("clinic_settings"),
		edge.From("internal_users", InternalUser.Type).Ref("internal_user_settings"),
		edge.From("patients", Patient.Type).Ref("patient_settings"),
		edge.To("clinic_customers", CustomerSettingOnClinics.Type),
	}
}

// Annotations of the Setting.
func (Setting) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "setting",
		},
	}
}
