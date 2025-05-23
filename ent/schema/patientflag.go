package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

// PatientFlag holds the schema definition for the PatientFlag entity.
type PatientFlag struct {
	ent.Schema
}

// Fields of the PatientFlag.
func (PatientFlag) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("patient_flag_id").StructTag(`json:"patient_flag_id"`),
		field.String("patient_flag_name").Unique(),
		field.String("patient_flag_display_name").Optional(),
		field.String("patient_flag_description").Optional(),
		field.Bool("patient_flag_is_active").Default(true),
		field.Time("patient_flag_created_at").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now),
		field.String("patient_flag_color").Default("#90EE90").Optional(),
		field.String("patient_flaged_by").Default("System"),
	}
}

// Edges of the PatientFlag.
func (PatientFlag) Edges() []ent.Edge {
	//TODO: Add link to patient after it's in
	return nil
}

func (PatientFlag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "patient_flag"},
	}
}
